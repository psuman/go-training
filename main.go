package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

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

	var svc service.ItemService

	cacheFinder := cache.Initialize("localhost:6379")
	dao := persistence.Initialize("mongodb://localhost:27017")
	svc = service.ItemCatalogService{CacheFinder: cacheFinder,
		ItemDao: dao}

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
	logger.Log("err", http.ListenAndServe(*listen, nil))
	fmt.Println("after service")
}
