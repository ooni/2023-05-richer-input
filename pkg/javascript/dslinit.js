"use strict"

ooni.dsl.compose = function(...args) {
	function compose2(left, right) {
		return {
			"stage_name": "compose",
			"arguments": {},
			"children": [left, right],
		}
	}

	function composeN(left, rights) {
		if (rights.length <= 0) {
			throw "composeN called with zero right functions"
		}
		if (rights.length == 1) {
			return compose2(left, rights[0])
		}
		return compose2(left, composeN(rights[0], rights.slice(1)))
	}

	if (args.length < 2) {
		throw "compose called with less that two functions"
	}
	return composeN(args[0], args.slice(1))
}

ooni.dsl.discard = function() {
	return {
		"stage_name": "discard",
		"arguments": {},
		"children": []
	}
}

ooni.dsl.dnsLookupGetaddrinfo = function() {
	return {
		"stage_name": "dns_lookup_getaddrinfo",
		"arguments": {},
		"children": []
	}
}

ooni.dsl.domainName = function(domain) {
	return {
		"stage_name": "domain_name",
		"arguments": {
			"domain": domain,
		},
		"children": []
	}
}

ooni.dsl.httpConnectionTls = function() {
	return {
		"stage_name": "http_connection_tls",
		"arguments": {},
		"children": []
	}
}

ooni.dsl.httpTransaction = function(options) {
	// TODO(bassosimone): implement all the options
	return {
		"stage_name": "http_transaction",
		"arguments": {},
		"children": []
	}
}

ooni.dsl.makeEndpointsForPort = function(port) {
	return {
		"stage_name": "make_endpoints_for_port",
		"arguments": {
			"port": port,
		},
		"children": []
	}
}

ooni.dsl.newEndpoint = function(address, options) {
	return {
		"stage_name": "new_endpoint",
		"arguments": {
			"endpoint": address,
			"domain": (options || {})["domain"] || "",
		},
		"children": []
	}
}

ooni.dsl.newEndpointPipeline = function(stage) {
	return {
		"stage_name": "new_endpoint_pipeline",
		"arguments": {},
		"children": [stage]
	}
}

ooni.dsl.tcpConnect = function() {
	return {
		"stage_name": "tcp_connect",
		"arguments": {},
		"children": []
	}
}

ooni.dsl.tlsHandshake = function(options) {
	return {
		"stage_name": "tls_handshake",
		"arguments": {
			"alpn": (options || {})["alpn"] || [],
			"skip_verify": (options || {})["skip_verify"] || false,
			"sni": (options || {})["sni"] || "",
			"x509_certs": (options || {})["x509_certs"] || [],
		},
		"children": []
	}
}
