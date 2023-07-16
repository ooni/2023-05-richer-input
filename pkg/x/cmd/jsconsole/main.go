package main

import (
	"os"

	"github.com/apex/log"
	"github.com/ooni/2023-05-richer-input/pkg/jsengine"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func main() {
	rtx := jsengine.New(log.Log)
	script := string(runtimex.Try1(os.ReadFile(os.Args[1])))
	runtimex.Try0(rtx.RunScript(os.Args[1], script))
}
