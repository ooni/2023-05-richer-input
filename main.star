load("ooni.star", "dsl")

pipeline = dsl.compose(
	dsl.domain_name("www.youtube.com"),
	dsl.dns_lookup_getaddrinfo(),
	dsl.make_endpoints_for_port(443),
	dsl.new_endpoint_pipeline(
		dsl.compose(
			dsl.tcp_connect(),
			dsl.tls_handshake(),
			dsl.http_connection_tls(),
			dsl.http_transaction(),
			dsl.discard(),
		)
	)
)

print(json.encode(pipeline))

results = dsl.run(pipeline)

#print(json.encode(results))
