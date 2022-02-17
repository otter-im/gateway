package handler

import (
	"github.com/go-oauth2/oauth2/v4/server"
	"net/http"
	"net/http/httputil"
	"sync"
)

type AuthProxyHandler struct {
	Server    *server.Server
	Host      string
	Scheme    string
	proxy     *httputil.ReverseProxy
	proxyOnce sync.Once
}

func (a *AuthProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.proxyOnce.Do(func() {
		a.proxy = &httputil.ReverseProxy{Director: director(a.Host, a.Scheme)}
	})
	token, err := a.Server.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Header.Set("X-Otter-Login-User-Id", token.GetUserID())
	a.proxy.ServeHTTP(w, r)
}

func director(host, scheme string) func(r *http.Request) {
	return func(r *http.Request) {
		r.URL.Host = host
		r.URL.Scheme = scheme
	}
}
