package httputils

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoHandle struct {
	Method  string
	Handler echo.HandlerFunc
}

type EchoRouter struct {
	port        string
	router      *echo.Echo
	prefix      string
	middleware  []echo.MiddlewareFunc
	handlers    map[string]EchoHandle
	healthCheck bool
	logRequest  bool
	cors        bool
	recovery    bool
}

func NewEchoRouter(port string) Router {
	return &EchoRouter{
		port:       port,
		router:     echo.New(),
		middleware: make([]echo.MiddlewareFunc, 0),
		handlers:   make(map[string]EchoHandle, 0),
	}
}

func (r *EchoRouter) Default() {
	r.
		AllowCors().
		AllowHealthCheck().
		AllowLog().
		AllowRecovery().
		ServeHTTP()
}

func (r *EchoRouter) AddPrefix(prefix string) Router {
	r.prefix = prefix
	return r
}

func (r *EchoRouter) AllowRecovery() Router {
	r.recovery = true
	return r
}

func (r *EchoRouter) AllowLog() Router {
	r.logRequest = true
	return r
}

func (r *EchoRouter) AllowHealthCheck() Router {
	r.healthCheck = true
	return r
}

func (r *EchoRouter) AllowCors() Router {
	r.cors = true
	return r
}

func (r *EchoRouter) ServeHTTP() {

	// middleware
	if r.logRequest {
		r.AddMiddleware(middleware.Logger())
	}
	if r.cors {
		r.AddMiddleware(middleware.CORS())
	}

	for _, h := range r.middleware {
		r.router.Use(h)
	}

	// handler
	if r.healthCheck {
		r.AddPath("/health", "GET", r.healthCheckHandler)
	}
	e := r.router.Group(r.prefix)

	for p, h := range r.handlers {
		e.Add(h.Method, p, h.Handler)
	}

	if r.recovery {
		e.Use(middleware.Recover())
	}

	// server
	go func() {
		log.Println("Server started on: " + r.port)
		if err := r.router.Start(":" + r.port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[CRITICAL] Shutting down the server: %+v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.router.Shutdown(ctx); err != nil {
		log.Fatalf("[CRITICAL] Server shutdown failed: %+v", err)
	}
}

func (EchoRouter) healthCheckHandler(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "service is running",
	})
}

func (r *EchoRouter) AddPath(path, method string, handler echo.HandlerFunc) *EchoRouter {
	r.handlers[path] = EchoHandle{
		Method:  method,
		Handler: handler,
	}
	return r
}

func (r *EchoRouter) AddMiddleware(middleware echo.MiddlewareFunc) *EchoRouter {
	r.middleware = append(r.middleware, middleware)
	return r
}
