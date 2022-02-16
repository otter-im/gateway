package handler

import (
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
	"github.com/google/uuid"
	"github.com/otter-im/gateway/internal/service"
	"github.com/otter-im/identity/pkg/rpc"
	"net/http"
	"os"
)

func LoginPageHandler(server *server.Server) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, err := session.Start(r.Context(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Method == http.MethodPost {
			if r.Form == nil {
				if err := r.ParseForm(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			response, err := service.LookupService().Authorize(r.Context(), &rpc.AuthorizationRequest{
				Username: r.Form.Get("username"),
				Password: r.Form.Get("password"),
			})

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			uid, err := uuid.FromBytes(response.GetId())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			store.Set("user_id", uid.String())
			store.Save()

			if r.Form.Get("prompt") == "none" {
				w.Header().Set("Location", "/oauth/authorize")
				w.WriteHeader(http.StatusFound)
			} else {
				w.Header().Set("Location", "/gateway")
				w.WriteHeader(http.StatusFound)
			}
			return
		}
		// TODO: Use the mux router instead of serving the page like this
		outputHtml(w, r, "web/static/login.html")
	}
}

func AuthPageHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("user_id"); !ok {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	// TODO: Use the mux router instead of serving the page like this
	outputHtml(w, r, "web/static/gateway.html")
}

func outputHtml(w http.ResponseWriter, r *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, r, file.Name(), fi.ModTime(), file)
}
