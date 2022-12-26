package httputils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type MuxHandle struct {
	Method  string
	Handler http.HandlerFunc
}

type MuxRouter struct {
	port        string
	router      *mux.Router
	prefix      string
	middleware  []func(next http.Handler) http.Handler
	handlers    map[string]MuxHandle
	healthCheck bool
	logRequest  bool
	cors        bool
	recovery    bool
}

func NewMuxRouter(port string) Router {
	return &MuxRouter{
		port:       port,
		router:     mux.NewRouter(),
		middleware: make([]func(next http.Handler) http.Handler, 0),
		handlers:   make(map[string]MuxHandle, 0),
	}
}

func (r *MuxRouter) Default() {
	r.
		AllowCors().
		AllowHealthCheck().
		AllowLog().
		AllowRecovery().
		ServeHTTP()
}

func (r *MuxRouter) AddPrefix(prefix string) Router {
	r.prefix = prefix
	return r
}

func (r *MuxRouter) AllowRecovery() Router {
	r.recovery = true
	return r
}

func (r *MuxRouter) AllowLog() Router {
	r.logRequest = true
	return r
}

func (r *MuxRouter) AllowHealthCheck() Router {
	r.healthCheck = true
	return r
}

func (r *MuxRouter) AllowCors() Router {
	r.cors = true
	return r
}

func (r *MuxRouter) ServeHTTP() {
	r.router.PathPrefix(r.prefix)
	// middleware
	if r.logRequest {
		r.AddMiddleware(r.logRequestlMiddleware)
	}
	if r.cors {
		r.AddMiddleware(r.accessControlMiddleware)
	}

	for _, h := range r.middleware {
		r.router.Use(h)
	}

	// handler
	if r.healthCheck {
		r.AddPath("/health", "GET", r.healthCheckHandler)
	}

	for p, h := range r.handlers {
		r.router.HandleFunc(p, h.Handler).Methods(h.Method)
	}
	if r.recovery {
		r.router.Use(r.recoverylMiddleware)
	}

	// server
	server := http.Server{
		Addr:         ":" + r.port,
		Handler:      r.router,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		log.Println("Server started on: " + r.port)
		log.Fatal(server.ListenAndServe())
	}()

	log.Println("exit", <-errs)
}

func (MuxRouter) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	MuxJson(w, 200, map[string]interface{}{
		"message": "service is running",
	})
}

func (MuxRouter) logRequestlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("request: " + r.RemoteAddr + ", method: " + r.Method + ", path: " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (MuxRouter) accessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (MuxRouter) recoverylMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				errorMessage := fmt.Sprintf("recovering from err %v with %s", err, buf)
				w.Write([]byte(`{"error":"+` + errorMessage + `"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (r *MuxRouter) AddPath(path, method string, handler http.HandlerFunc) *MuxRouter {
	r.handlers[path] = MuxHandle{
		Method:  method,
		Handler: handler,
	}
	return r
}

func (r *MuxRouter) AddMiddleware(middleware func(next http.Handler) http.Handler) *MuxRouter {
	r.middleware = append(r.middleware, middleware)
	return r
}

func MuxJson(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
