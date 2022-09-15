package httputils

type Router interface {
	Default()
	ServeHTTP()
	AddPrefix(prefix string) Router
	AllowRecovery() Router
	AllowLog() Router
	AllowHealthCheck() Router
	AllowCors() Router
}
