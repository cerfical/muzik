package webapp

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/log"
)

func New(args []string) *App {
	var a App

	a.Log = log.New()
	a.Config = loadConfig(args, a.Log)
	a.Log = a.Log.WithLevel(a.Config.Log.Level)

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
	a.Log.WithFields("addr", a.Config.Server.Addr).Info("Starting up the server")

	serv := http.Server{
		Addr:    a.Config.Server.Addr,
		Handler: setupMiddleware(setupRouter(a.routes), a.middleware),
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.Log.WithError(err).Error("Server terminated abnormally")
		}

		// Make sure the outer goroutine unblocks in case the server was terminated before any signals arrived
		signal.Stop(sigChan)
		close(sigChan)
	}()

	// If the server was terminated due to some other reason, return immediately
	if _, ok := <-sigChan; !ok {
		return
	}
	a.Log.Info("Shutting down the server")

	// Try to shutdown the server cleanly and if that fails, close the server
	if err := serv.Shutdown(context.Background()); err != nil {
		a.Log.WithError(err).Error("Server shutdown failed")
		serv.Close()
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
