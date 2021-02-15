package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
)

// Gateway default service
type Gateway struct {
	Name       string
	Endpoint   string
	EnableHTTP bool
	ser        *http.Server
}

// NewGateway create gateway
func NewGateway(name, endpoint string, enableHTTP bool) *Gateway {
	return &Gateway{
		Name:       name,
		Endpoint:   endpoint,
		EnableHTTP: enableHTTP,
	}
}

// GetName get service name
func (ser *Gateway) GetName() string {
	return ser.Name
}

// Serve start service
func (ser *Gateway) Serve(opt ...grpc.ServerOption) error {
	if ser.ser != nil {
		return errors.New("service: service is already serve")
	}
	mux := NewMux(ser.EnableHTTP, opt...)
	routeInfos := mux.GetRouteInfos()
	for _, info := range routeInfos {
		bytes, _ := json.Marshal(info)
		fmt.Println(bytes)
	}
	ser.ser = &http.Server{
		Addr:    ser.Endpoint,
		Handler: mux,
	}

	return ser.ser.ListenAndServe()
}

// Close close service
func (ser *Gateway) Close() error {
	return ser.ser.Close()
}
