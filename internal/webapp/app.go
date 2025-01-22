package webapp

import (
	"net/http"

	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/httpserv"
	"github.com/cerfical/muzik/internal/log"
	"github.com/cerfical/muzik/internal/middleware"
)

func New(args []string) *App {
	var a App

	a.Log = log.New()
	a.Config = loadConfig(args, a.Log)
	a.Log = a.Log.WithLevel(a.Config.Log.Level)
	a.Server = &httpserv.Server{Addr: a.Config.Server.Addr}

	a.Use(middleware.LogRequest(a.Log))

	return &a
}

func loadConfig(args []string, l *log.Logger) *config.Config {
	c, err := config.Load(args)
	if err != nil {
		l.WithError(err).Fatal("Error loading config file")
	}
	return c
}

type App struct {
	Config *config.Config
	Server *httpserv.Server
	Log    *log.Logger

	middleware []Middleware
	routes     []route
}

type Middleware func(http.Handler) http.Handler

type route struct {
	path    string
	handler http.HandlerFunc
}

func (a *App) Use(m Middleware) {
	a.middleware = append(a.middleware, m)
}

func (a *App) Route(path string, h http.HandlerFunc) {
	a.routes = append(a.routes, route{path, h})
}

func (a *App) Run() {
	a.Log.WithFields("serv.addr", a.Server.Addr).Info("Starting up the server")

	a.Server.Handler = setupMiddleware(setupRouter(a.routes), a.middleware)
	if err := a.Server.Run(); err != nil {
		a.Log.WithError(err).Fatal("Server terminated abnormally")
	}

	a.Log.Info("Shutting down the server")
	if err := a.Server.Close(); err != nil {
		a.Log.WithError(err).Fatal("Server shutdown failed")
	}
}

func setupRouter(routes []route) http.Handler {
	mux := http.NewServeMux()
	for _, r := range routes {
		mux.HandleFunc(r.path, r.handler)
	}
	return mux
}

func setupMiddleware(h http.Handler, m []Middleware) http.Handler {
	for _, m := range m {
		h = m(h)
	}
	return h
}
