package main

import (
	"fmt"
	"os"

	"github.com/ooni/2023-05-richer-input/pkg/javascript"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func main() {
	rtx := runtimex.Try1(javascript.NewRuntime())
	script := string(runtimex.Try1(os.ReadFile(os.Args[1])))
	value := runtimex.Try1(rtx.RunScript("", script))
	fmt.Printf("%+v\n", value)
}
