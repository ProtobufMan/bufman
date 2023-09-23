package main

import (
	"fmt"
	"github.com/ProtobufMan/bufman/internal/config"
	"github.com/ProtobufMan/bufman/internal/dal"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/router"
	"github.com/ProtobufMan/bufman/internal/task"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

func main() {
	config.LoadConfig()

	model.InitDB()

	dal.SetDefault(config.DataBase)

	// init router
	r := router.InitRouter()

	// init async task
	err := task.InitAsyncTask()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = task.AsyncTaskManager.Stop()
	}()

	err = http.ListenAndServe(
		fmt.Sprintf(":%v", config.Properties.BufMan.Port),
		// For gRPC clients, it's convenient to support HTTP/2 without TLS. You can
		// avoid x/net/http2 by using http.ListenAndServeTLS.
		h2c.NewHandler(r, &http2.Server{}),
	)
	if err != nil {
		panic(err)
	}
}
