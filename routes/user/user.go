package user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jtejido/go-github/config"
	"github.com/jtejido/go-github/core"
	"github.com/jtejido/go-github/service"
	"net/http"
	"sort"
)

func userHandler(conf *config.Config, svc *service.Github) func(c *gin.Context) {
	return func(c *gin.Context) {
		ch := make(chan *service.PublicUser)
		errChan := make(chan error)

		req := c.Request
		req.ParseForm()
		r := req.Form
		val, ok := r["name"]
		if !ok {
			c.JSON(http.StatusBadRequest, core.NewErrorResponseWithCode("name required", 400))
			return
		} else {
			if len(val) > conf.MaxLimit {
				c.JSON(http.StatusBadRequest, core.NewErrorResponseWithCode(fmt.Sprintf("the limit for accepted has been reached: %v", conf.MaxLimit), 400))
				return
			}
			sort.Strings(val)
		}

		resp := make([]*service.PublicUser, 0)
		for _, v := range val {
			go func(v string) {
				user, err := svc.GetUser(v)
				if err != nil {
					errChan <- err
				}
				ch <- user
			}(v)
		}

		for i := 0; i < len(val); i++ {
			select {
			case err := <-errChan:
				c.JSON(http.StatusInternalServerError, core.NewErrorResponseWithCode(err.Error(), 500))
				return
			case item := <-ch:
				if item != nil {
					resp = append(resp, item)
				}
			}
		}

		close(ch)
		close(errChan)

		c.JSON(http.StatusOK, resp)
	}
}

func Setup(ctx context.Context, conf *config.Config, svc *service.Github, g *gin.RouterGroup) {
	g.GET("", userHandler(conf, svc))
}
