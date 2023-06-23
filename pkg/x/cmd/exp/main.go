package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/ooni/2023-05-richer-input/pkg/x/dsl"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func main() {
	template := dsl.Compose(
		dsl.String("www.example.com"),
		dsl.Getaddrinfo(),
		dsl.MakeEndpointList(443),
		dsl.MakeEndpointPipeline(
			dsl.QUICHandshake(
			/*
				dsl.QUICHandshakeOptionALPN("h3"),
				dsl.QUICHandshakeOptionSNI("www.example.com"),
				dsl.QUICHandshakeOptionSkipVerify(true),
			*/
			),
		),
	)

	{
		data := runtimex.Try1(json.Marshal(template))
		fmt.Fprintf(os.Stderr, "%s\n", string(data))
	}

	registry := dsl.NewFunctionRegistry()
	function, err := dsl.CompileInvocation(registry, template)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	rtx := dsl.NewRuntime(dsl.RuntimeOptionLogger(log.Log))
	defer rtx.Close()

	result := function.Apply(context.Background(), rtx, &dsl.Void{})
	log.Infof("%T %+v", result, result)

	observations := rtx.ExtractObservations()
	fmt.Printf("%s\n", string(runtimex.Try1(json.Marshal(observations))))
}
