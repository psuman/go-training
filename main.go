package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	persistence "github.com/psuman/go-training/service/persistence"

	cache "github.com/psuman/go-training/service/cache"

	service "github.com/psuman/go-training/service"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	var (
		listen = flag.String("listen", ":9090", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var svc service.ItemService

	cacheFinder := cache.RedisCacheFinder{}
	cacheFinder = cacheFinder.Initialize("localhost:6379", logger)

	dao := persistence.MongoItemDao{}
	dao = dao.Initialize("mongodb://localhost:27017", logger)

	svc = service.ItemCatalogService{CacheFinder: cacheFinder,
		ItemDao: dao, Logger: logger}

	findItemHandler := httptransport.NewServer(
		service.MakeFindItemEndPoint(svc),
		service.DecodeFindItemRequest,
		service.EncodeResponse,
	)

	addItemHandler := httptransport.NewServer(
		service.MakeAddItemEndPoint(svc),
		service.DecodeAddItemRequest,
		service.EncodeResponse,
	)

	fmt.Println("inside service")
	http.Handle("/find-item", findItemHandler)
	http.Handle("/add-item", addItemHandler)
	// http.Handle("/metrics", promhttp.Handler())

	dao.FindAll()
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

	defer func() {
		svc.Close()
	}()

}
