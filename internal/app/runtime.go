package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/gorilla/mux"
	"github.com/otter-im/gateway/internal/app/handler"
	"github.com/otter-im/gateway/internal/config"
	"github.com/otter-im/gateway/internal/service"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc/connectivity"
	"log"
	mathRand "math/rand"
	"net"
	"net/http"
	"time"
)

var (
	exitHooks = make([]func() error, 0)
)

func Init() error {
	log.Printf("environment: %s\n", config.ServiceEnvironment())

	rand.Seed(uint64(time.Now().UnixNano()))
	mathRand.Seed(time.Now().UnixNano())
	AddExitHook(service.ExitHook)

	if err := checkPostgres(); err != nil {
		return err
	}

	if err := checkRedis(); err != nil {
		return err
	}

	if err := checkIdentity(); err != nil {
		return err
	}
	return nil
}

func Run() error {
	router := mux.NewRouter()

	tokenStore := &OtterTokenStore{}
	clientStore := &OtterClientStore{}
	srv := initServer(router, tokenStore, clientStore)

	router.HandleFunc("/login", handler.LoginPageHandler())
	router.HandleFunc("/auth", handler.AuthPageHandler)

	router.Handle("/profile", &handler.AuthProxyHandler{Server: srv, Host: config.APIProfileHost(), Scheme: config.APIProfileScheme()})
	router.Handle("/profile/{id}", &handler.AuthProxyHandler{Server: srv, Host: config.APIProfileHost(), Scheme: config.APIProfileScheme()})

	http.Handle("/", router)
	addr := net.JoinHostPort(config.ServiceHost(), config.ServicePort())
	log.Printf("listening on: %v\n", addr)
	return http.ListenAndServe(addr, nil)
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
	srv.CheckCodeChallengeMethod(oauth2.CodeChallengeS256)
	srv.SetUserAuthorizationHandler(handler.UserAuthorizationHandler)

	srv.SetInternalErrorHandler(handler.InternalErrorHandler)
	srv.SetResponseErrorHandler(handler.ResponseErrorHandler)

	router.HandleFunc("/oauth/authorize", handler.AuthorizeHandler(srv))
	router.HandleFunc("/oauth/token", handler.TokenHandler(srv))
	return srv
}

func checkPostgres() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := Postgres().Ping(ctx); err != nil {
		return fmt.Errorf("postgresql connection failure: %v", err)
	}
	return nil
}

func checkRedis() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if cmd := RedisRing().Ping(ctx); cmd.Err() != nil {
		return fmt.Errorf("redis connection failure: %v", cmd.Err())
	}
	return nil
}

func checkIdentity() error {
	conn := service.IdentityConn()
	if conn == nil {
		return errors.New("previous failed connection")
	}
	if state := conn.GetState(); state == connectivity.TransientFailure || state == connectivity.Shutdown {
		return errors.New("gRPC is in an invalid state")
	}
	return nil
}
