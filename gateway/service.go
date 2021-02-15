package gateway

import (
	"google.golang.org/grpc"
)

// Service gateway
type Service interface {
	// GetName get service name
	GetName() string
	// Serve start service. parameter endpoint string, enableHttp bool
	Serve(...grpc.ServerOption) error
	// Close close service
	Close() error
}

// Run service
func Run(service Service, opt ...grpc.ServerOption) error {
	defer service.Close()

	return service.Serve(opt...)
}
