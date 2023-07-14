package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net/http"

	"github.com/apex/log"
	"github.com/ooni/probe-engine/pkg/netxlite"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

// apiMustFetchCA fetches, validates and returns the CA or PANICS.
func apiMustFetchCA() string {
	log.Info("- fetching the CA")

	resp := runtimex.Try1(http.Get("https://black.riseup.net/ca.crt"))
	runtimex.Assert(resp.StatusCode == 200, "unexpected HTTP response status")
	defer resp.Body.Close()

	log.Infof("HTTP response: %+v", resp)

	body := string(runtimex.Try1(netxlite.ReadAllContext(context.Background(), resp.Body)))
	log.Infof("fetched CA:\n%s\n", string(body))
	return body
}

// apiMustFetchEIPService fetches and parses the [*apiEIPService] or PANICS.
func apiMustFetchEIPService(caCert string) *apiEIPService {
	log.Info("- fetching eip-service.json")

	// create and fill a certificate pool
	pool := x509.NewCertPool()
	runtimex.Assert(pool.AppendCertsFromPEM([]byte(caCert)), "AppendCertsFromPEM failed")

	// create a client using a transport using the pool
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: pool,
			},
		},
	}

	// perform the HTTP round trip
	resp := runtimex.Try1(client.Get("https://api.black.riseup.net/3/config/eip-service.json"))
	runtimex.Assert(resp.StatusCode == 200, "unexpected HTTP response status")
	defer resp.Body.Close()

	log.Infof("HTTP response: %+v", resp)

	// read the whole body
	body := runtimex.Try1(netxlite.ReadAllContext(context.Background(), resp.Body))
	log.Infof("fetched eip-service.json:\n%s\n", string(body))

	// parse the response body
	var eipService apiEIPService
	runtimex.Try0(json.Unmarshal(body, &eipService))
	return &eipService
}
