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

	external_invoker "github.com/psuman/go-training/service/external_invoker"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		listen = flag.String("listen", ":9090", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger

	logFile, err := os.Create("app.log")

	if err != nil {
		panic(err)
	}

	logger = log.NewLogfmtLogger(logFile)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	fieldKeys := []string{"method", "error"}

	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "item_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "item_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc service.ItemService

	cacheFinder := cache.RedisCacheFinder{}
	cacheFinder = cacheFinder.Initialize("localhost:6379", logger)

	dao := persistence.MongoItemDao{}
	dao = dao.Initialize("mongodb://localhost:27017", logger)

	httpClient := &http.Client{}

	extSvc := external_invoker.ExternalFindItemServiceInvokerImpl{ServiceUrl: "http://localhost:7070/find-ext-item", Timeout: 10, HttpClient: httpClient}

	svc = service.ItemCatalogService{CacheFinder: cacheFinder,
		ItemDao: dao, ExtService: extSvc, Logger: logger}

	svc = service.MetricsMiddleware{requestCount, requestLatency, svc}

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

	http.Handle("/find-item", findItemHandler)
	http.Handle("/add-item", addItemHandler)
	http.Handle("/metrics", promhttp.Handler())

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
