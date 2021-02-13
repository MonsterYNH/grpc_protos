package route

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

// DealGQLGateway deal gql gateway
func DealGQLGateway(filePrefix string, plugin *protogen.Plugin, descs []*descriptorpb.FileDescriptorProto) error {
	infos, err := parseGQLProtoFile(descs)
	if err != nil {
		return err
	}

	return generateGQLFile(infos, filePrefix, plugin)
}

func parseGQLProtoFile(descs []*descriptorpb.FileDescriptorProto) ([]Info, error) {
	// We're dealing with each descriptors to out wrapper struct
	// in order to access easily plugin options, package name, comment, etc...
	var files []*spec.File
	for _, f := range descs {
		files = append(files, spec.NewFile(f, false))
	}

	infos := make([]Info, 0)
	for _, file := range files {
		for _, service := range file.Services() {

			log.Println(service.Name(), "======")
			for _, method := range service.Methods() {
				infos = append(infos, Info{
					ServiceName:  service.Name(),
					Method:       method.Schema.GetType().String(),
					Pattern:      method.Schema.GetName(),
					RequestType:  method.Input(),
					ResponseType: method.Output(),
				})
				log.Println(method.Schema.GetType(), "======")
				log.Println(method.Schema.GetName(), "======")
				log.Println(method.Input(), "======")
				log.Println(method.Output(), "======")
			}
		}
	}
	return infos, nil
}

func generateGQLFile(Infos []Info, filePrefix string, plugin *protogen.Plugin) error {
	var buf bytes.Buffer

	tmpl, err := template.New("gql").Parse(gqlGatewayTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(&buf, Infos); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s%s.go", filePrefix, ".gql.gateway")
	template := plugin.NewGeneratedFile(fileName, ".")

	template.Write(buf.Bytes())

	return nil
}
