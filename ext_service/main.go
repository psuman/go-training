package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	ext_service "github.com/psuman/go-training/ext_service/service"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	var (
		listen = flag.String("listen", ":7070", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger

	logFile, err := os.Create("ext_service.log")

	if err != nil {
		panic(err)
	}

	logger = log.NewLogfmtLogger(logFile)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var svc ext_service.ItemService

	svc = ext_service.ItemCatalogService{Logger: logger}

	findItemHandler := httptransport.NewServer(
		ext_service.MakeFindItemEndPoint(svc),
		ext_service.DecodeFindItemRequest,
		ext_service.EncodeResponse,
	)

	http.Handle("/find-ext-item", findItemHandler)

	errs := make(chan error, 2)

	go func() {
		logger.Log("transport", "http", "address", *listen, "msg", "listening")
		errs <- http.ListenAndServe(*listen, nil)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}
