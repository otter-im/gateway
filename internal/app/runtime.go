package app

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/gorilla/mux"
	"github.com/otter-im/auth/internal/config"
	"github.com/otter-im/auth/internal/handler"
	"golang.org/x/exp/rand"
	mathRand "math/rand"
	"net"
	"net/http"
	"time"
)

var (
	exitHooks = make([]func() error, 0)
)

func Init() {
	rand.Seed(uint64(time.Now().UnixNano()))
	mathRand.Seed(time.Now().UnixNano())
}

func Run() error {
	router := mux.NewRouter()

	tokenStore := &OtterTokenStore{}
	clientStore := &OtterClientStore{}
	srv := initServer(router, tokenStore, clientStore)

	router.HandleFunc("/login", handler.LoginPageHandler)
	router.HandleFunc("/auth", handler.AuthPageHandler)
	router.HandleFunc("/test", handler.TestHandler(srv))

	http.Handle("/", router)
	return http.ListenAndServe(net.JoinHostPort(config.ServiceHost(), config.ServicePort()), nil)
}

func Exit() error {
	for _, hook := range exitHooks {
		err := hook()
		if err != nil {
			return err
		}
	}
	return nil
}

func AddExitHook(hook func() error) {
	exitHooks = append(exitHooks, hook)
}

func initServer(router *mux.Router, tokenStore oauth2.TokenStore, clientStore oauth2.ClientStore) *server.Server {
	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)
	srv.SetUserAuthorizationHandler(handler.UserAuthorizationHandler)

	srv.SetInternalErrorHandler(handler.InternalErrorHandler)
	srv.SetResponseErrorHandler(handler.ResponseErrorHandler)

	router.HandleFunc("/oauth/authorize", handler.AuthorizeHandler(srv))
	router.HandleFunc("/oauth/token", handler.TokenHandler(srv))
	return srv
}
