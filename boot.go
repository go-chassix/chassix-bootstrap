package bootstrap

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"c5x.io/bootstrap/config"
)

//StartHttpServer starting a http server for restful api
func StartHttpServer(handler http.Handler) {
	fmt.Println(strconv.Itoa(config.Server().Data.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Server().Data.Port), handler))
}

//StartRPCServer starting a gRPC server
func StartGrpcServer() {

}

func init() {
	config.WatchBootstrapConfig()
	config.Load()
}
