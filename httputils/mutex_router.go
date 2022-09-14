package httputils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

type PathHandle struct {
	Method  string
	Handler http.HandlerFunc
}

type MuxRouter struct {
	port        string
	router      *mux.Router
	prefix      string
	middleware  []func(next http.Handler) http.Handler
	pathHandle  map[string]PathHandle
	healthCheck bool
	logRequest  bool
	cors        bool
	recovery    bool
}

func NewMuxRouter(port string) *MuxRouter {
	return &MuxRouter{
		port:       port,
		router:     mux.NewRouter(),
		middleware: make([]func(next http.Handler) http.Handler, 0),
		pathHandle: make(map[string]PathHandle, 0),
	}
}

func (r *MuxRouter) AddPrefix(prefix string) *MuxRouter {
	r.prefix = prefix
	return r
}

func (r *MuxRouter) AllowLog() *MuxRouter {
	r.logRequest = true
	return r
}

func (r *MuxRouter) AllowHealthCheck() *MuxRouter {
	r.healthCheck = true
	return r
}

func (r *MuxRouter) AllowCors() *MuxRouter {
	r.cors = true
	return r
}

func (r *MuxRouter) ServeHTTP() {
	r.router.PathPrefix(r.prefix)
	// middleware
	if r.logRequest {
		r.AddMiddleware(logRequest)
	}
	if r.cors {
		r.AddMiddleware(accessControlMiddleware)
	}

	for _, h := range r.middleware {
		r.router.Use(h)
	}

	// handler
	if r.healthCheck {
		r.AddPath("/health", "GET", healthCheck)
	}

	for p, h := range r.pathHandle {
		r.router.HandleFunc(p, h.Handler).Methods(h.Method)
	}
	if r.recovery {
		r.router.Use(recovery)
	}

	// server
	server := http.Server{
		Addr:         ":" + r.port,
		Handler:      r.router,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	log.Println("Server started on: " + r.port)
	log.Fatal(server.ListenAndServe())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	Json(w, 200, map[string]interface{}{
		"message": "service is running",
	})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("request: " + r.RemoteAddr + ", method: " + r.Method + ", path: " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func accessControlMiddleware(next http.Handler) http.Handler {
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

func recovery(next http.Handler) http.Handler {
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
	r.pathHandle[path] = PathHandle{
		Method:  method,
		Handler: handler,
	}
	return r
}

func (r *MuxRouter) AddMiddleware(middleware func(next http.Handler) http.Handler) *MuxRouter {
	r.middleware = append(r.middleware, middleware)
	return r
}

func Json(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
