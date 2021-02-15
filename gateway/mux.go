package gateway

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MonsterYNH/protoc-gen-gateway/route"

	"google.golang.org/grpc"
)

// Mux route
type Mux struct {
	grpcHandler  *grpc.Server
	httpHandlers map[string]http.Handler
	EnableHTTP   bool
}

// NewMux create mux
func NewMux(enableHTTP bool, opt ...grpc.ServerOption) *Mux {
	return &Mux{
		grpcHandler:  grpc.NewServer(opt...),
		httpHandlers: make(map[string]http.Handler),
		EnableHTTP:   enableHTTP,
	}
}

func (mux *Mux) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	scheme := req.Header.Get("Content-Type")
	switch scheme {
	case "application/grpc":
		mux.grpcHandler.ServeHTTP(resp, req)
	case "application/json":
		if !mux.EnableHTTP {
			resp.WriteHeader(http.StatusNotFound)
			resp.Write([]byte(`{message: "not enable http service"}`))
		}
		// route by root path
		paths := strings.Split(req.URL.Path, "/")
		if len(paths) == 0 {
			resp.WriteHeader(http.StatusNotFound)
			resp.Write([]byte(`{message: "more information please visit /doc"}`))
			return
		}
		handler, exist := mux.httpHandlers[paths[0]]
		if exist {
			handler.ServeHTTP(resp, req)
			return
		}
	}
	resp.WriteHeader(http.StatusNotFound)
	resp.Write([]byte(`{message: "bad request"}`))
}

// RegistGRPC add grpc
func (mux *Mux) RegistGRPC(ser *grpc.ServiceDesc, impl interface{}) {
	mux.grpcHandler.RegisterService(ser, impl)
}

// RegistHandler add http handler
func (mux *Mux) RegistHandler(name string, handler http.Handler) error {
	if _, exist := mux.httpHandlers[name]; exist {
		return fmt.Errorf("gateway: http handler %s is already exist", name)
	}
	mux.httpHandlers[name] = handler
	return nil
}

// GetRouteInfos get route infos
func (mux *Mux) GetRouteInfos() []route.Info {
	infos := route.NewInfos()
	// get grpc route infos
	for names, serviceInfo := range mux.grpcHandler.GetServiceInfo() {
		name := strings.Split(names, ".")
		for _, service := range serviceInfo.Methods {
			if err := infos.RegistInfo(route.Info{
				ServiceName: name[0],
				Type:        "grpc",
				Pattern:     fmt.Sprintf("%s/%s", names, service.Name),
			}); err != nil {
				panic(err)
			}
		}
	}

	if mux.EnableHTTP {
		if err := infos.RegistInfos(route.HTTPInfos.Infos); err != nil {
			panic(err)
		}
	}

	return infos.Infos
}
