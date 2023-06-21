package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/apex/log"
	"github.com/ooni/2023-05-richer-input/pkg/x/dslx"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func evaluate(ctx context.Context, env *dslx.Environment) *dslx.Observations {
	// create a measurement function
	function := dslx.Compose(
		env.DomainToResolve("example.com"),
		env.DNSLookupParallel(
			env.DNSLookupGetaddrinfo(),
			env.DNSLookupUDP(dslx.DNSLookupUDPOptionEndpoint("[2001:4860:4860::8844]:53")),
			env.DNSLookupUDP(dslx.DNSLookupUDPOptionEndpoint("8.8.4.4:53")),
		),
		env.EndpointParallel(
			env.EndpointPipeline(
				dslx.EndpointPort(443),
				env.TCPConnect(),
				env.TLSHandshake(),
				env.Discard(),
			),
			env.EndpointPipeline(
				dslx.EndpointPort(80),
				env.TCPConnect(),
				env.Discard(),
			),
			env.EndpointPipeline(
				dslx.EndpointPort(443),
				env.QUICHandshake(),
				env.Discard(),
			),
		),
	)

	// evaluate the measurement function
	return env.Eval(ctx, function)
}

func main() {
	// create a context
	ctx := context.Background()

	// create the execution environment
	env := dslx.NewEnvironment(&atomic.Int64{}, log.Log, time.Now())

	// be careful with execution
	observations, err := env.Try(func() *dslx.Observations {
		return evaluate(ctx, env)
	})

	// handle any panic that may have occurred
	if err != nil {
		log.Fatalf("FATAL: %s", err.Error())
	}

	// serialize and print observations
	fmt.Printf("%s\n", string(runtimex.Try1(json.Marshal(observations))))
}
