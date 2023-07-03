package dsl_test

//
// QA tests to make sure we handle error conditions correctly in terms of
// typing and of generated results when running RunnableASTNode.
//
// We have a fixed testing scenario with www.example.com and www.example.org
// and we generate several prossible measurement conditions.
//

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/apex/log"
	"github.com/google/gopacket/layers"
	"github.com/ooni/2023-05-richer-input/pkg/dsl"
	"github.com/ooni/netem"
	"github.com/ooni/probe-engine/pkg/netemx"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// qaWebServerAddress is the address of www.example.com and www.example.org.
const qaWebServerAddress = "93.184.216.34"

//
// This section of the file contains functions to generate AST nodes
//

func qaPipelineDNS() dsl.Stage[string, *dsl.DNSLookupResult] {
	return dsl.DNSLookupParallel(
		dsl.DNSLookupGetaddrinfo(),
		dsl.DNSLookupUDP(net.JoinHostPort(netemx.DefaultServersResolver, "53")),
	)
}

func qaPipelineHTTP() dsl.Stage[*dsl.DNSLookupResult, *dsl.Void] {
	return dsl.Compose(
		dsl.MakeEndpointsForPort(80),
		dsl.NewEndpointPipeline(
			dsl.Compose4(
				dsl.TCPConnect(),
				dsl.HTTPConnectionTCP(),
				dsl.HTTPTransaction(),
				dsl.Discard[*dsl.HTTPResponse](),
			),
		),
	)
}

func qaPipelineHTTPS() dsl.Stage[*dsl.DNSLookupResult, *dsl.Void] {
	return dsl.Compose(
		dsl.MakeEndpointsForPort(443),
		dsl.NewEndpointPipeline(
			dsl.Compose5(
				dsl.TCPConnect(),
				dsl.TLSHandshake(),
				dsl.HTTPConnectionTLS(),
				dsl.HTTPTransaction(),
				dsl.Discard[*dsl.HTTPResponse](),
			),
		),
	)
}

func qaPipelineHTTP3() dsl.Stage[*dsl.DNSLookupResult, *dsl.Void] {
	return dsl.Compose(
		dsl.MakeEndpointsForPort(443),
		dsl.NewEndpointPipeline(
			dsl.Compose4(
				dsl.QUICHandshake(),
				dsl.HTTPConnectionQUIC(),
				dsl.HTTPTransaction(),
				dsl.Discard[*dsl.HTTPResponse](),
			),
		),
	)
}

func qaNewMeasurementPipelineForDomain(domain string) dsl.Stage[*dsl.Void, *dsl.Void] {
	return dsl.Compose3(
		dsl.DomainName(domain),
		qaPipelineDNS(),
		dsl.MeasureMultipleEndpoints(
			qaPipelineHTTP(),
			qaPipelineHTTPS(),
			qaPipelineHTTP3(),
		),
	)
}

func qaNewRunnableASTNode() dsl.RunnableASTNode {
	pipeline := dsl.RunStagesInParallel(
		qaNewMeasurementPipelineForDomain("www.example.com"),
		qaNewMeasurementPipelineForDomain("www.example.org"),
	)
	ast := runtimex.Try1(json.Marshal(pipeline.ASTNode()))
	var loadable dsl.LoadableASTNode
	runtimex.Try0(json.Unmarshal(ast, &loadable))
	loader := dsl.NewASTLoader()
	return runtimex.Try1(loader.Load(&loadable))
}

//
// This section of the file contains code to generate environments
//

func qaNewEnvironment(clientDNSConfig *netem.DNSConfig) *netemx.Environment {
	// clientConfig configures the client topology
	clientConfig := &netemx.ClientConfig{
		ClientAddr:   "", // use the default
		DNSConfig:    clientDNSConfig,
		ResolverAddr: "", // use the default
	}

	// create the configuration of the uncensored DNS server.
	serversDNSConfig := netem.NewDNSConfig()
	serversDNSConfig.AddRecord("www.example.com", "www.example.com", qaWebServerAddress)
	serversDNSConfig.AddRecord("www.example.org", "www.example.org", qaWebServerAddress)

	// create the overall configuration for the servers.
	serversConfig := &netemx.ServersConfig{
		DNSConfig:    serversDNSConfig,
		ResolverAddr: "", // use the default
		Servers: []netemx.ConfigServerStack{{
			ServerAddr: qaWebServerAddress,
			HTTPServers: []netemx.ConfigHTTPServer{{
				Port:    80,
				QUIC:    false,
				Handler: nil, // use the default
			}, {
				Port:    443,
				QUIC:    false,
				Handler: nil, // use the default
			}, {
				Port:    443,
				QUIC:    true,
				Handler: nil, // use the default
			}},
		}},
	}

	// return the environment
	return netemx.NewEnvironment(clientConfig, serversConfig)
}

//
// This section of the file runs the measurement pipeline
//

func qaRunNode(runnable dsl.RunnableASTNode) (*dsl.Observations, error) {
	input := dsl.NewValue(&dsl.Void{}).AsGeneric()
	rtx := dsl.NewMeasurexliteRuntime(log.Log, time.Now())
	if err := dsl.Try(runnable.Run(context.Background(), rtx, input)); err != nil {
		return nil, err
	}
	return dsl.ReduceObservations(rtx.ExtractObservations()...), nil

}

//
// This section of the file contains tests
//

func TestQASuccess(t *testing.T) {
	dnsConfig := netem.NewDNSConfig()
	dnsConfig.AddRecord("www.example.com", "www.example.com", qaWebServerAddress)
	dnsConfig.AddRecord("www.example.org", "www.example.org", qaWebServerAddress)

	env := qaNewEnvironment(dnsConfig)
	defer env.Close()

	var (
		observations *dsl.Observations
		err          error
	)

	env.Do(func() {
		observations, err = qaRunNode(qaNewRunnableASTNode())
	})

	if err != nil {
		t.Fatal(err)
	}

	// TODO(bassosimone): check the observations
	_ = observations
}

func TestQADNSLookupGetaddrinfoFailure(t *testing.T) {
	dnsConfig := netem.NewDNSConfig()
	// Note: we're not filling the DNS config, which causes NXDOMAIN

	env := qaNewEnvironment(dnsConfig)
	defer env.Close()

	var (
		observations *dsl.Observations
		err          error
	)

	env.Do(func() {
		observations, err = qaRunNode(qaNewRunnableASTNode())
	})

	if err != nil {
		t.Fatal(err)
	}

	// TODO(bassosimone): check the observations
	_ = observations
}

func TestQADNSLookupUDPFailure(t *testing.T) {
	dnsConfig := netem.NewDNSConfig()
	dnsConfig.AddRecord("www.example.com", "www.example.com", qaWebServerAddress)
	dnsConfig.AddRecord("www.example.org", "www.example.org", qaWebServerAddress)

	env := qaNewEnvironment(dnsConfig)
	defer env.Close()

	// Note: this rule should prevent UDP communication
	env.DPIEngine().AddRule(&netem.DPIDropTrafficForServerEndpoint{
		Logger:          log.Log,
		ServerIPAddress: netemx.DefaultServersResolver,
		ServerPort:      53,
		ServerProtocol:  layers.IPProtocolUDP,
	})

	var (
		observations *dsl.Observations
		err          error
	)

	env.Do(func() {
		observations, err = qaRunNode(qaNewRunnableASTNode())
	})

	if err != nil {
		t.Fatal(err)
	}

	// TODO(bassosimone): check the observations
	_ = observations
}

func TestQAFullDNSSpoofing(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test in short mode")
	}

	dnsConfig := netem.NewDNSConfig()
	dnsConfig.AddRecord("www.example.com", "www.example.com", qaWebServerAddress)
	dnsConfig.AddRecord("www.example.org", "www.example.org", qaWebServerAddress)

	env := qaNewEnvironment(dnsConfig)
	defer env.Close()

	// Note: this rule should spoof the responses
	env.DPIEngine().AddRule(&netem.DPISpoofDNSResponse{
		Addresses: []string{"10.10.34.35"},
		Logger:    log.Log,
		Domain:    "www.example.com",
	})
	env.DPIEngine().AddRule(&netem.DPISpoofDNSResponse{
		Addresses: []string{"10.10.34.35"},
		Logger:    log.Log,
		Domain:    "www.example.org",
	})

	var (
		observations *dsl.Observations
		err          error
	)

	env.Do(func() {
		observations, err = qaRunNode(qaNewRunnableASTNode())
	})

	if err != nil {
		t.Fatal(err)
	}

	// TODO(bassosimone): check the observations
	_ = observations
}

func TestQATCPConnectFailure(t *testing.T) {
	t.Run("in case of failure connecting on 443/tcp", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skip test in short mode")
		}

		dnsConfig := netem.NewDNSConfig()
		dnsConfig.AddRecord("www.example.com", "www.example.com", qaWebServerAddress)
		dnsConfig.AddRecord("www.example.org", "www.example.org", qaWebServerAddress)

		env := qaNewEnvironment(dnsConfig)
		defer env.Close()

		// Note: this rule should prevent connecting
		env.DPIEngine().AddRule(&netem.DPIDropTrafficForServerEndpoint{
			Logger:          log.Log,
			ServerIPAddress: qaWebServerAddress,
			ServerPort:      443,
			ServerProtocol:  layers.IPProtocolTCP,
		})

		var (
			observations *dsl.Observations
			err          error
		)

		env.Do(func() {
			observations, err = qaRunNode(qaNewRunnableASTNode())
		})

		if err != nil {
			t.Fatal(err)
		}

		// TODO(bassosimone): check the observations
		_ = observations
	})

	t.Run("in case of failure connecting on 80/tcp", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skip test in short mode")
		}

		dnsConfig := netem.NewDNSConfig()
		dnsConfig.AddRecord("www.example.com", "www.example.com", qaWebServerAddress)
		dnsConfig.AddRecord("www.example.org", "www.example.org", qaWebServerAddress)

		env := qaNewEnvironment(dnsConfig)
		defer env.Close()

		// Note: this rule should prevent connecting
		env.DPIEngine().AddRule(&netem.DPIDropTrafficForServerEndpoint{
			Logger:          log.Log,
			ServerIPAddress: qaWebServerAddress,
			ServerPort:      80,
			ServerProtocol:  layers.IPProtocolTCP,
		})

		var (
			observations *dsl.Observations
			err          error
		)

		env.Do(func() {
			observations, err = qaRunNode(qaNewRunnableASTNode())
		})

		if err != nil {
			t.Fatal(err)
		}

		// TODO(bassosimone): check the observations
		_ = observations
	})
}

func TestQATLSHandshakeFailure(t *testing.T) {
	dnsConfig := netem.NewDNSConfig()
	dnsConfig.AddRecord("www.example.com", "www.example.com", qaWebServerAddress)
	dnsConfig.AddRecord("www.example.org", "www.example.org", qaWebServerAddress)

	env := qaNewEnvironment(dnsConfig)
	defer env.Close()

	// Note: this rule should reset the handshake
	env.DPIEngine().AddRule(&netem.DPIResetTrafficForTLSSNI{
		Logger: log.Log,
		SNI:    "www.example.com",
	})

	var (
		observations *dsl.Observations
		err          error
	)

	env.Do(func() {
		observations, err = qaRunNode(qaNewRunnableASTNode())
	})

	if err != nil {
		t.Fatal(err)
	}

	// TODO(bassosimone): check the observations
	_ = observations
}

func TestQAQUICHandshakeFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("skip test in short mode")
	}

	dnsConfig := netem.NewDNSConfig()
	dnsConfig.AddRecord("www.example.com", "www.example.com", qaWebServerAddress)
	dnsConfig.AddRecord("www.example.org", "www.example.org", qaWebServerAddress)

	env := qaNewEnvironment(dnsConfig)
	defer env.Close()

	// Note: this rule should prevent handshaking
	env.DPIEngine().AddRule(&netem.DPIDropTrafficForServerEndpoint{
		Logger:          log.Log,
		ServerIPAddress: qaWebServerAddress,
		ServerPort:      443,
		ServerProtocol:  layers.IPProtocolUDP,
	})

	var (
		observations *dsl.Observations
		err          error
	)

	env.Do(func() {
		observations, err = qaRunNode(qaNewRunnableASTNode())
	})

	if err != nil {
		t.Fatal(err)
	}

	// TODO(bassosimone): check the observations
	_ = observations
}
