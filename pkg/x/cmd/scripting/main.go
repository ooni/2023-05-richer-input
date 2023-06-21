package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/apex/log"
	"github.com/ooni/2023-05-richer-input/pkg/x/dslx"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func evaluate(ctx context.Context, env *dslx.Environment, expession []any) *dslx.Observations {
	// create a measurement function
	function := env.Compile(expession)

	// evaluate the measurement function
	return env.Eval(ctx, function)
}

func main() {
	// create a context
	ctx := context.Background()

	// create the execution environment
	env := dslx.NewEnvironment(&atomic.Int64{}, log.Log, time.Now())

	// read the input JSON
	data := runtimex.Try1(os.ReadFile(os.Args[1]))
	var expression []any
	runtimex.Try0(json.Unmarshal(data, &expression))

	// be careful with execution
	observations, err := env.Try(func() *dslx.Observations {
		return evaluate(ctx, env, expression)
	})

	// handle any panic that may have occurred
	if err != nil {
		log.Fatalf("FATAL: %s", err.Error())
	}

	// serialize and print observations
	fmt.Printf("%s\n", string(runtimex.Try1(json.Marshal(observations))))
}
