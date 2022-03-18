package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

func New() *Github {
	return &Github{
		apiUrl:    apiUrl,
		client:    &http.Client{},
		RateLimit: new(RateLimit),
	}
}

func (g *Github) GetUser(username string) (*PublicUser, error) {
	url := g.apiUrl + strings.Replace(user_url, ":user", username, -1)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := g.client.Do(req)
	defer resp.Body.Close()

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

	return user, nil
}

func (g *Github) ListUsers(results int) []*PublicUser {
	url := g.apiUrl + users_url + "?per_page=" + fmt.Sprintf("%d", results)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := g.client.Do(req)
	defer resp.Body.Close()

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
		panic(err)
	}
	users := make([]*PublicUser, results)
	err = json.Unmarshal(contents, &users)
	if err != nil {
		panic(err)
	}

	return users
}
