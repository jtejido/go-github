package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jtejido/go-github/cache"
	"github.com/jtejido/go-github/config"
	// "log"
	"net/http"
	"time"
)

type logWriter struct {
	gin.ResponseWriter
	cache *cache.Cache
	conf  *config.Config
}

func (w bodyCacheWriter) Write(b []byte) (int, error) {
	status := w.Status()
	var res map[string]interface{}
	json.Unmarshal(b, &res)
	if status == http.StatusOK {
		if v, ok := res["login"]; ok {
			// uncomment to see that it writes in cache upon success (200)
			// log.Printf("Writing body for Login: %s in cache.", v.(string))
			expire := time.Now().Add(time.Duration(w.conf.UserLifetime) * time.Second)
			w.cache.Set(v.(string), &cache.Item{b, &expire})
		}
	}
	return w.ResponseWriter.Write(b)
}

func CacheCheck(conf *config.Config, store *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("name")
		item, err := store.Get(id)
		if err == nil {
			// uncomment to see that it fetches from cache, if present, it returns the cached
			// item, otherwise, it will continue and proceed with the github API url
			// log.Println("item found.")
			c.Data(http.StatusOK, "application/json", item.Value)
			c.Abort()
		} else {
			c.Writer = bodyCacheWriter{cache: store, conf: conf, ResponseWriter: c.Writer}
			c.Next()
		}
	}
}
