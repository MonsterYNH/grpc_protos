package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/MonsterYNH/protoc-gen-gateway/route"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// Protoc passes pluginpb.CodeGeneratorRequest in via stdin
	// marshalled with Protobuf
	input, _ := ioutil.ReadAll(os.Stdin)
	var req pluginpb.CodeGeneratorRequest
	proto.Unmarshal(input, &req)

	// Initialise our plugin with default options
	opts := protogen.Options{}
	plugin, err := opts.New(&req)
	if err != nil {
		panic(err)
	}

	file := plugin.Files[len(plugin.Files)-1]

	if err := route.DealHTTPRouteInfo(file, plugin); err != nil {
		panic(err)
	}

	if err := route.DealGQLRouteInfo(file, plugin, req.GetProtoFile()); err != nil {
		panic(err)
	}

	// Generate a response from our plugin and marshall as protobuf
	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	if err != nil {
		panic(err)
	}

	// Write the response to stdout, to be picked up by protoc
	fmt.Fprintf(os.Stdout, string(out))
}
