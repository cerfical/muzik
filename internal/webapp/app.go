package webapp

import (
	stdlog "log"

	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cerfical/muzik/internal/config"
	"github.com/cerfical/muzik/internal/log"

	"github.com/gorilla/mux"
)

func New(args []string) *App {
	var a App

	a.Log = log.New()
	a.Config = loadConfig(args, a.Log)
	a.Log = a.Log.WithLevel(a.Config.Log.Level)
	a.router = mux.NewRouter()

	a.Use(fillPathParams)

	return &a
}

// fillPathParams makes gorilla/mux path parameters available via [http.Request.PathValue].
func fillPathParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range mux.Vars(r) {
			r.SetPathValue(key, val)
		}

		next.ServeHTTP(w, r)
	})
}

func loadConfig(args []string, l *log.Logger) *config.Config {
	c, err := config.Load(args)
	if err != nil {
		l.Fatal("error loading config file", err)
	}
	return c
}

type App struct {
	Config *config.Config
	Log    *log.Logger

	NotFound         http.HandlerFunc
	MethodNotAllowed http.HandlerFunc

	router *mux.Router
}

type Middleware func(http.Handler) http.Handler

func (a *App) Use(m Middleware) {
	a.router.Use(mux.MiddlewareFunc(m))
}

func (a *App) Route(method, path string, h http.HandlerFunc) {
	a.router.HandleFunc(path, h).Methods(method)
}

func (a *App) Run() {
	a.Log.WithFields("addr", a.Config.Server.Addr).Info("starting up the server")

	if a.MethodNotAllowed != nil {
		a.router.MethodNotAllowedHandler = a.MethodNotAllowed
	}

	if a.NotFound != nil {
		a.router.NotFoundHandler = a.NotFound
	}

	serv := http.Server{
		Addr:     a.Config.Server.Addr,
		Handler:  a.router,
		ErrorLog: stdlog.New(&httpErrorLog{a.Log}, "", 0),
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := serv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.Log.Error("server terminated abnormally", err)
		}

		// Make sure the outer goroutine unblocks in case the server was terminated before any signals arrived
		signal.Stop(sigChan)
		close(sigChan)
	}()

	// If the server was terminated due to some other reason, return immediately
	if _, ok := <-sigChan; !ok {
		return
	}
	a.Log.Info("shutting down the server")

	// Try to shutdown the server cleanly and if that fails, close the server
	if err := serv.Shutdown(context.Background()); err != nil {
		a.Log.Error("server shutdown failed", err)
		serv.Close()
	}
}

type httpErrorLog struct {
	*log.Logger
}

func (w *httpErrorLog) Write(p []byte) (int, error) {
	// Trim carriage return produced by stdlog
	n := len(p)
	if n > 0 && p[n-1] == '\n' {
		p = p[0 : n-1]
		n--
	}

	w.Error("HTTP serve error", errors.New(string(p)))
	return n, nil
}
