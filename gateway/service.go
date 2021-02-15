package gateway

import (
	"fmt"

	"github.com/MonsterYNH/protoc-gen-gateway/route"
)

// Service gateway
type Service interface {
	// GetName get service name
	GetName() string
	// GetEndpoint
	GetEndpoint() string
	// GetRouteInfos get route infos
	GetRouteInfos() []route.Info
	// Serve start service. parameter endpoint string, enableHttp bool
	Serve() error
	// Close close service
	Close() error
}

// Run service
func Run(service Service) error {
	defer service.Close()

	fmt.Println(fmt.Sprintf("service %s start at %s", service.GetName(), service.GetEndpoint()))

	routeInfos := service.GetRouteInfos()
	for _, info := range routeInfos {
		fmt.Println(fmt.Sprintf("route: %s %s %s %s %s %s", info.ServiceName, info.Type, info.Method, info.Pattern, info.RequestType, info.ResponseType))
	}

	return service.Serve()
}
