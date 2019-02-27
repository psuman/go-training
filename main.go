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

	var svc service.FindItemService
	svc = service.FindItemInCatalogService{CacheFinder: cache.RedisCacheFinder{},
		ItemDao: persistence.MongoItemDao{}}

	findItemHandler := httptransport.NewServer(
		service.MakeFindItemEndPoint(svc),
		service.DecodeFindItemRequest,
		service.EncodeResponse,
	)

	fmt.Println("inside service")
	http.Handle("/find-item", findItemHandler)
	// http.Handle("/metrics", promhttp.Handler())
	logger.Log("err", http.ListenAndServe(*listen, nil))
	fmt.Println("after service")
}
