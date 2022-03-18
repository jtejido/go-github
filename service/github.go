package service

import (
	"encoding/json"
	"github.com/jtejido/go-github/cache"
	"github.com/jtejido/go-github/config"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	apiUrl     = "https://api.github.com"
	user_url   = "/users/:user"
	users_url  = "/users"
	dateLayout = "2006-01-02T15:04:05Z"
)

type Github struct {
	apiUrl    string
	client    *http.Client
	RateLimit *RateLimit
	cache     *cache.Cache
	conf      *config.Config
}

type RateLimit struct {
	Limit     int64 `json:"limit"`
	Remaining int64 `json:"remaining"`
	Reset     int64 `json:"reset"`
}

type SimpleUser struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

type PublicUser struct {
	SimpleUser
	Name        string `json:"name"`
	Company     string `json:"company"`
	Blog        string `json:"blog"`
	Location    string `json:"location"`
	Email       string `json:"email"`
	Bio         string `json:"bio"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
}

func New(cache *cache.Cache, conf *config.Config) *Github {
	return &Github{
		apiUrl:    apiUrl,
		client:    &http.Client{},
		RateLimit: new(RateLimit),
		cache:     cache,
		conf:      conf,
	}
}

func (g *Github) GetUser(username string) (*PublicUser, error) {
	if item, err := g.cache.Get(username); err == nil {
		user := new(PublicUser)
		if err := json.Unmarshal(item.Value, &user); err != nil {
			return nil, err
		}

		return user, nil
	}

	url := g.apiUrl + strings.Replace(user_url, ":user", username, -1)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}

	limit, err := strconv.ParseInt(resp.Header.Get("X-RateLimit-Limit"), 10, 64)
	if err == nil {
		g.RateLimit.Limit = limit
	}
	remaining, err := strconv.ParseInt(resp.Header.Get("X-RateLimit-Remaining"), 10, 64)
	if err == nil {
		g.RateLimit.Remaining = remaining
	}
	reset, err := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 64)
	if err == nil {
		g.RateLimit.Reset = reset
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	user := new(PublicUser)
	err = json.Unmarshal(contents, &user)
	if err != nil {
		return nil, err
	}

	expire := time.Now().Add(time.Duration(g.conf.UserLifetime) * time.Second)
	g.cache.Set(username, &cache.Item{contents, &expire})
	// think about using ETags for stale users in the cache

	return user, nil
}
