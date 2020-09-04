package chassix

import (
	"c5x.io/chassix/internal"
	"fmt"
	"net/http"
	"strconv"

	"c5x.io/logx"
)

const (
	ModuleApp            = internal.KeyAppConfig
	ModuleApollo         = internal.KeyApolloConfig
	ModuleServer         = internal.KeyServerConfig
	ModuleLogging        = internal.KeyLoggingConfig
	ModuleDatasource     = internal.KeyDatasourceConfig
	ModuleSqlxDatasource = internal.KeySqlxDatasourceConfig
	ModuleRedis          = internal.KeyRedisConfig
	ModuleCache          = internal.KeyCacheConfig
)

//StartHttpServer starting a http server for restful api
func StartHttpServer(handler http.Handler) {
	var log = logx.New().Category("boot").Component("starter")

	log.Infof("Server starting... IP: %s, port: %d", internal.Server().Data.Addr, internal.Server().Data.Port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(internal.Server().Data.Port), handler))
}

//StartRPCServer starting a gRPC server
func StartGrpcServer() {
	//todo x
}

//Init 在所有模块注册就绪后初始化
func Init() {
	internal.Load()
	logx.SetConfig(internal.Logging())
}

func init() {
	fmt.Print(welcome3)
	internal.WatchBootstrapConfig()
}

type Module struct {
	Name      string
	ConfigPtr interface{}
}

func Register(module *Module) error {
	internal.Watch(module.Name, module.ConfigPtr)
	return nil
}
