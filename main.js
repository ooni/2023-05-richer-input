"use strict"

const ooni = require("./ooni.js")

const pipeline = ooni.dsl.compose(
	ooni.dsl.domainName("www.youtube.com"),
	ooni.dsl.dnsLookupGetaddrinfo(),
	ooni.dsl.makeEndpointsForPort(443),
	ooni.dsl.newEndpointPipeline(
		ooni.dsl.compose(
			ooni.dsl.tcpConnect(),
			ooni.dsl.tlsHandshake(),
			ooni.dsl.httpConnectionTls(),
			ooni.dsl.httpTransaction(),
			ooni.dsl.discard(),
		)
	)
)

console.log(JSON.stringify(pipeline))

const results = ooni.dsl.run(pipeline)

console.log(JSON.stringify(results))
