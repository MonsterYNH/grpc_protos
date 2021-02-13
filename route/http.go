package route

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
	options "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// DealHTTPGateway deal http gateway
func DealHTTPGateway(file *protogen.File, plugin *protogen.Plugin) error {
	routeInfos, err := parseHTTPProtoFile(file)
	if err != nil {
		return err
	}

	return generateHTTPFile(routeInfos, file.GeneratedFilenamePrefix, plugin)
}

// parseHTTPProtoFile parse http proto
func parseHTTPProtoFile(file *protogen.File) ([]Info, error) {
	f := &descriptor.File{
		FileDescriptorProto: file.Proto,
		GoPkg: descriptor.GoPackage{
			Path: string(file.GoImportPath),
			Name: string(file.GoPackageName),
		},
	}

	routes := make([]Info, 0)
	for _, service := range f.GetService() {
		for _, method := range service.GetMethod() {
			log.Println(method.GetName())
			serviceName := method.GetName()
			var httpMethod, httpPattern string
			ext := proto.GetExtension(method.Options, options.E_Http)
			httpRule := ext.(*options.HttpRule)
			if len(httpRule.GetGet()) != 0 {
				httpMethod = "get"
				httpPattern = httpRule.GetGet()
			} else if len(httpRule.GetPatch()) != 0 {
				httpMethod = "patch"
				httpPattern = httpRule.GetPatch()
			} else if len(httpRule.GetPost()) != 0 {
				httpMethod = "post"
				httpPattern = httpRule.GetPost()
			} else if len(httpRule.GetDelete()) != 0 {
				httpMethod = "delete"
				httpPattern = httpRule.GetDelete()
			} else if len(httpRule.GetPut()) != 0 {
				httpMethod = "put"
				httpPattern = httpRule.GetPut()
			} else {
				return nil, errors.New("route: no match http method")
			}
			requestType := method.GetInputType()
			responseType := method.GetOutputType()
			routes = append(routes, Info{
				ServiceName:  serviceName,
				Method:       httpMethod,
				Pattern:      httpPattern,
				RequestType:  requestType,
				ResponseType: responseType,
			})
		}
	}

	return routes, nil
}

// generateHTTPFile generate http file
func generateHTTPFile(Infos []Info, filePrefix string, plugin *protogen.Plugin) error {
	var buf bytes.Buffer

	tmpl, err := template.New("http").Parse(httpGatewayTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(&buf, Infos); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s%s.go", filePrefix, ".http.gateway")
	template := plugin.NewGeneratedFile(fileName, ".")

	template.Write(buf.Bytes())

	return nil
}
