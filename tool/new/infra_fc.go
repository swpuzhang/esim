package new

func init() {
	Files = append(Files, infrafc1, infrafc2, infrafc3)
}

var (
	infrafc1 = &FileContent{
		FileName: "infra.go",
		Dir:      "internal/infra",
		Content: `package infra

import (
	"sync"
	"github.com/google/wire"
	"github.com/jukylin/esim/container"
	"github.com/jukylin/esim/mysql"
	"github.com/jukylin/esim/grpc"
	"{{.ProPath}}{{.ServerName}}/internal/infra/repo"
)

//Do not change the function name and var name
//  infraOnce
//  onceInfra
//  infraSet
//  NewInfra

var infraOnce sync.Once
var onceInfra *Infra

type Infra struct {
	*container.Esim

	DB *mysql.Client

	GrpcClient *grpc.Client

	UserRepo repo.UserRepo
}

var infraSet = wire.NewSet(
	wire.Struct(new(Infra), "*"),
	provideDb,
	provideUserRepo,
)


func NewInfra() *Infra {
	infraOnce.Do(func() {
		esim  := container.NewEsim()
		onceInfra = initInfra(esim, provideGrpcClient(esim))
	})

	return onceInfra
}

func NewStubsInfra(grpcClient *grpc.Client) *Infra {
	infraOnce.Do(func() {
		esim  := container.NewEsim()
		onceInfra = initInfra(esim, grpcClient)
	})

	return onceInfra
}

// Close close the infra when app stop
func (this *Infra) Close()  {

	this.DB.Close()
}

func (this *Infra) HealthCheck() []error {
	var errs []error
	var err error

	dbErrs := this.DB.Ping()
	if err != nil{
		errs = append(errs, dbErrs...)
	}

	return errs
}


func provideDb(esim *container.Esim) *mysql.Client {

	clientOptions := mysql.ClientOptions{}
	mysqlClent := mysql.NewClient(
		clientOptions.WithConf(esim.Conf),
		clientOptions.WithLogger(esim.Logger),
		clientOptions.WithProxy(
			func() interface{} {
				monitorProxyOptions := mysql.MonitorProxyOptions{}
				return mysql.NewMonitorProxy(
					monitorProxyOptions.WithLogger(esim.Logger),
					monitorProxyOptions.WithConf(esim.Conf),
					monitorProxyOptions.WithTracer(esim.Tracer),
				)
			},
		),
	)

	return mysqlClent
}


func provideUserRepo(esim *container.Esim) repo.UserRepo {
	return repo.NewDBUserRepo(esim.Logger)
}


func provideGrpcClient(esim *container.Esim) *grpc.Client {
	clientOptional := grpc.ClientOptionals{}
	clientOptions := grpc.NewClientOptions(
		clientOptional.WithLogger(esim.Logger),
		clientOptional.WithConf(esim.Conf),
	)

	grpcClient := grpc.NewClient(clientOptions)

	return grpcClient
}

`,
	}

	infrafc2 = &FileContent{
		FileName: "wire.go",
		Dir:      "internal/infra",
		Content: `//+build wireinject

package infra

import (
	"github.com/google/wire"
	"github.com/jukylin/esim/grpc"
	"github.com/jukylin/esim/container"
)


func initInfra(esim *container.Esim,grpc *grpc.Client) *Infra {
	wire.Build(infraSet)
	return nil
}
`,
	}

	infrafc3 = &FileContent{
		FileName: "wire_gen.go",
		Dir:      "internal/infra",
		Content: `// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package infra

import (
	"github.com/jukylin/esim/container"
	"github.com/jukylin/esim/grpc"
)

// Injectors from wire.go:

func initInfra(esim *container.Esim, grpc2 *grpc.Client) *Infra {
	mysqlClient := provideDb(esim)
	userRepo := provideUserRepo(esim)
	infra := &Infra{
		Esim:     esim,
		DB:       mysqlClient,
		GrpcClient: grpc2,
		UserRepo: userRepo,
	}
	return infra
}
`,
	}
)
