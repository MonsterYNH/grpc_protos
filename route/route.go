package route

import (
	"fmt"

	"google.golang.org/grpc"
)

// Info route info
type Info struct {
	ServiceName  string
	Type         string
	Method       string
	Pattern      string
	RequestType  string
	ResponseType string
}

const (
	// RequestTypeHTTP http request
	RequestTypeHTTP = "http"
	// RequestTypeGRPC grpc request
	RequestTypeGRPC = "grpc"
	// RequestTypeGQL gql request
	RequestTypeGQL = "gql"
)

var (
	// HTTPInfos infos
	HTTPInfos = NewInfos()
)

// Infos route infos
type Infos struct {
	Infos      []Info
	patternMap map[string]struct{}
}

// NewInfos create infos
func NewInfos() Infos {
	return Infos{
		Infos:      make([]Info, 0),
		patternMap: make(map[string]struct{}),
	}
}

// RegistInfo regist info
func (infos *Infos) RegistInfo(info Info) error {
	if info.Type != RequestTypeHTTP && info.Type != RequestTypeGRPC && info.Type == RequestTypeGQL {
		return fmt.Errorf("route: route type %s is not support", info.Type)
	}
	if _, exist := infos.patternMap[info.Pattern]; exist {
		return fmt.Errorf("route: pattern %s is already exist", info.Pattern)
	}
	infos.Infos = append(infos.Infos, info)
	return nil
}

// RegistInfos regist infos
func (infos *Infos) RegistInfos(data []Info) error {
	for _, info := range data {
		if err := infos.RegistInfo(info); err != nil {
			return err
		}
	}
	return nil
}

// GetInfos get route infos
func (infos *Infos) GetInfos() []Info {
	return infos.Infos
}

// GetGrpcInfos get grpc infos
func GetGrpcInfos(ser *grpc.Server) []Info {
	for name, info := range ser.GetServiceInfo() {
		for _, method := range info.Methods {
			fmt.Println(name, method.Name)
		}
	}
	return nil
}
