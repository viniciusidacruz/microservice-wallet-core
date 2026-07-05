package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	return &WebServer{
		Router:        router,
		WebServerPort: serverPort,
	}
}

func (w *WebServer) Get(path string, handler http.HandlerFunc) {
	w.Router.Get(path, handler)
}

func (w *WebServer) Post(path string, handler http.HandlerFunc) {
	w.Router.Post(path, handler)
}

func (w *WebServer) Delete(path string, handler http.HandlerFunc) {
	w.Router.Delete(path, handler)
}

func (w *WebServer) Start() {
	http.ListenAndServe(w.WebServerPort, w.Router)
}
