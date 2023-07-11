package main

import (
	"github.com/ProtobufMan/bufman/internal/config"
	"github.com/ProtobufMan/bufman/internal/router"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

func main() {
	config.LoadConfig()

	// init router
	r := router.InitRouter()

	err := http.ListenAndServe(
		"localhost:39099",
		// For gRPC clients, it's convenient to support HTTP/2 without TLS. You can
		// avoid x/net/http2 by using http.ListenAndServeTLS.
		h2c.NewHandler(r, &http2.Server{}),
	)
	panic(err)
}
