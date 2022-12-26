package main

import "github.com/huylqbk/codesample/httputils"

func main() {
	httputils.
		NewEchoRouter("8080").
		AddPrefix("/v1").
		AllowCors().
		AllowRecovery().
		AllowHealthCheck().
		ServeHTTP()
}
