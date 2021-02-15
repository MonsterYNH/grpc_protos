package gateway

import (
	"net/http"

	"github.com/MonsterYNH/protoc-gen-gateway/route"
	"google.golang.org/grpc"
)

// Gateway default service
type Gateway struct {
	Name       string
	Endpoint   string
	EnableHTTP bool
	mux        *Mux
	ser        *http.Server
}

// NewGateway create gateway
func NewGateway(name, endpoint string, enableHTTP bool, opt ...grpc.ServerOption) *Gateway {
	return &Gateway{
		Name:       name,
		Endpoint:   endpoint,
		EnableHTTP: enableHTTP,
		mux:        NewMux(enableHTTP, opt...),
		ser:        &http.Server{},
	}
}

// GetName get service name
func (ser *Gateway) GetName() string {
	return ser.Name
}

// GetRouteInfos get route infos
func (ser *Gateway) GetRouteInfos() []route.Info {
	return ser.mux.GetRouteInfos()
}

// GetEndpoint get endpoint
func (ser *Gateway) GetEndpoint() string {
	return ser.Endpoint
}

// Serve start service
func (ser *Gateway) Serve() error {
	ser.ser = &http.Server{
		Addr:    ser.Endpoint,
		Handler: ser.mux,
	}

	return ser.ser.ListenAndServe()
}

// Close close service
func (ser *Gateway) Close() error {
	return ser.ser.Close()
}
