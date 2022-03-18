package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jtejido/go-github/cache"
	"github.com/jtejido/go-github/config"
	"github.com/jtejido/go-github/core"
	"github.com/jtejido/go-github/service"
	"net/http"
	"sort"
	"time"
)

func userHandler(svc *service.Github, conf *config.Config, store *cache.Cache) func(c *gin.Context) {
	return func(c *gin.Context) {
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

		resp := make([]*service.PublicUser, len(val))
		for k, v := range val {
			if item, err := store.Get(v); err == nil {
				user := new(service.PublicUser)
				if err := json.Unmarshal(item.Value, &user); err != nil {
					c.JSON(http.StatusInternalServerError, core.NewErrorResponseWithCode(err.Error(), 500))
					return
				}
				resp[k] = user
			} else {
				user, err := svc.GetUser(v)
				if err != nil {
					c.JSON(http.StatusInternalServerError, core.NewErrorResponseWithCode(err.Error(), 500))
					return
				}
				resp[k] = user
				expire := time.Now().Add(time.Duration(conf.UserLifetime) * time.Second)
				reqBodyBytes := new(bytes.Buffer)
				json.NewEncoder(reqBodyBytes).Encode(user)
				store.Set(v, &cache.Item{reqBodyBytes.Bytes(), &expire})
			}

		}

		c.JSON(http.StatusOK, resp)
	}
}

func Setup(ctx context.Context, conf *config.Config, svc *service.Github, c *cache.Cache, g *gin.RouterGroup) {
	g.GET("", userHandler(svc, conf, c))
}
