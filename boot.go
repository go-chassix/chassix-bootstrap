package chassix

import (
	"fmt"
	"net/http"
	"strconv"

	"c5x.io/chassix/config"
	"c5x.io/logx"
)

//StartHttpServer starting a http server for restful api
func StartHttpServer(handler http.Handler) {
	var log = logx.New().Category("boot").Component("starter")

	log.Infof("Server starting... IP: %s, port: %d", config.Server().Data.Addr, config.Server().Data.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Server().Data.Port), handler))
}

func LoadConfig() {
	config.Load()
}

//StartRPCServer starting a gRPC server
func StartGrpcServer() {
	//todo x
}

//Init 在所有模块注册就绪后初始化
func Init() {
	config.Load()
	logx.SetConfig(config.Logging())
}

func init() {
	fmt.Print(welcome3)
	config.WatchBootstrapConfig()
}

type Module struct {
	Name      string
	ConfigPtr interface{}
}

func Register(module *Module) error {
	config.Watch(module.Name, module.ConfigPtr)
	return nil
}
