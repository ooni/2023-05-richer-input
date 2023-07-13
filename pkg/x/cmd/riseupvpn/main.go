package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"time"

	"github.com/apex/log"
	"github.com/ooni/2023-05-richer-input/pkg/dsl"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// isTCPGatewayAccessible returns whether a gateways is accessible.
func isTCPGatewayAccessible(stage dsl.Stage[*dsl.Void, *dsl.Void]) bool {
	metrics := dsl.NewAccountingMetrics()
	rtx := dsl.NewMeasurexliteRuntime(log.Log, metrics, time.Now())
	input := dsl.NewValue(&dsl.Void{})
	ctx := context.Background()
	runtimex.Try0(dsl.Try(stage.Run(ctx, rtx, input)))
	expect := map[string]int64{
		"tcp_connect_success_count": 1,
	}
	return reflect.DeepEqual(expect, metrics.Snapshot())
}

// generateGatewaysDSL generates a DSL to measure each gateway listed by eipService.
func generateGatewaysDSL(eipService *apiEIPService) (output []dsl.Stage[*dsl.Void, *dsl.Void]) {
	for _, gw := range eipService.Gateways {
		for _, txp := range gw.Capabilities.Transport {
			if !txp.typeIsOneOf("obfs4", "openvpn") {
				continue
			}
			if !txp.supportsTCP() {
				continue
			}
			for _, port := range txp.Ports {
				stage := dslRuleMeasureGatewayReachability(gw.IPAddress, port)
				epnt := net.JoinHostPort(gw.IPAddress, port)
				log.Infof("- checking whether %s/tcp is accessible", epnt)
				if !isTCPGatewayAccessible(stage) {
					log.Warnf("gateway %s/tcp IS NOT ACCESSIBLE", epnt)
					continue
				}
				log.Infof("gateway %s/tcp is accessible", epnt)
				output = append(output, stage)
			}
		}
	}
	return
}

// mustGenerateDSL generates the DSL for measuring riseupvpn or PANICS on failure.
//
// The returned DSL roughly includes:
//
// - a stage to fetch the CA required by riseupvpn;
//
// - a stage to fetch provider.json;
//
// - a stage to fetch eip-service.json;
//
// - a stage to query the riseupvpn geo service;
//
// - a bunch of stages containing gateways to measure.
//
// The stages will run in parallel with reasonably small parallelism.
func mustGenerateDSL(eipService *apiEIPService, rootCA string) dsl.Stage[*dsl.Void, *dsl.Void] {
	// start with a list of stages generated from the reachable gateways
	stages := generateGatewaysDSL(eipService)

	// add a stage to measure fetching the CA file
	stages = append(stages, dslRuleFetchCA())

	// add a stage to measure fetching the provider.json file
	stages = append(stages, dslRuleFetchProviderURL())

	// add a stage to measure fetching the eip-services.json file
	stages = append(stages, dslRuleFetchEIPServiceURL(rootCA))

	// add a stage to measure fetching the geo services file
	stages = append(stages, dslRuleFetchGeoServiceURL(rootCA))

	// return the composed pipeline
	return dsl.RunStagesInParallel(stages...)
}

func main() {
	// fetch the CA file
	rootCA := apiMustFetchCA()

	// fetch the EIP services file
	eipServices := apiMustFetchEIPService(rootCA)

	// generate the DSL
	DSL := mustGenerateDSL(eipServices, rootCA)

	// serialize to JSON
	rawDSL := runtimex.Try1(json.Marshal(DSL.ASTNode()))

	// load the DSL
	var loadable dsl.LoadableASTNode
	runtimex.Try0(json.Unmarshal(rawDSL, &loadable))

	// make the DSL runnable
	loader := dsl.NewASTLoader()
	runnable := runtimex.Try1(loader.Load(&loadable))

	// make sure we can run the DSL
	log.Info("- checking whether we can run the generated DSL")
	ctx := context.Background()
	input := dsl.NewValue(&dsl.Void{}).AsGeneric()
	rtx := dsl.NewMinimalRuntime(log.Log)
	runtimex.Try0(dsl.Try(runnable.Run(ctx, rtx, input)))

	// dump the raw DSL on the stdout
	fmt.Printf("%s\n", rawDSL)
}
