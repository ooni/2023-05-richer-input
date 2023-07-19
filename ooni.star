def _compose(*stages):
    # make sure we have at least one stage
    if len(stages) < 1:
        fail("compose: expected at least one stage")

    # step 1: invert the list
    inverted = []
    for stage in stages:
        inverted.insert(0, stage)

    # step 2: construct the composition
    composed = None
    for stage in inverted:
        if composed == None:
            composed = stage
        else:
            composed = {
                "stage_name": "compose",
                "arguments": {},
                "children": [stage, composed],
            }
    return composed


def _discard():
    return {"stage_name": "discard", "arguments": {}, "children": []}


def _dns_lookup_getaddrinfo():
    return {"stage_name": "dns_lookup_getaddrinfo", "arguments": {}, "children": []}


def _domain_name(domain):
    return {
        "stage_name": "domain_name",
        "arguments": {
            "domain": domain,
        },
        "children": [],
    }


def _http_connection_tls():
    return {"stage_name": "http_connection_tls", "arguments": {}, "children": []}


def _http_transaction(**options):
    # TODO(bassosimone): implement all the options
    return {"stage_name": "http_transaction", "arguments": {}, "children": []}


def _make_endpoints_for_port(port):
    return {
        "stage_name": "make_endpoints_for_port",
        "arguments": {
            "port": port,
        },
        "children": [],
    }


def _new_endpoint(address, **options):
    return {
        "stage_name": "new_endpoint",
        "arguments": {
            "endpoint": address,
            "domain": options.get("domain", ""),
        },
        "children": [],
    }


def _new_endpoint_pipeline(stage):
    return {"stage_name": "new_endpoint_pipeline", "arguments": {}, "children": [stage]}


def _run(pipeline):
    return json.decode(_ooni.run_dsl(json.encode(pipeline)))


def _tcp_connect():
    return {"stage_name": "tcp_connect", "arguments": {}, "children": []}


def _tls_handshake(**options):
    return {
        "stage_name": "tls_handshake",
        "arguments": {
            "alpn": options.get("alpn", []),
            "skip_verify": options.get("skip_verify", False),
            "sni": options.get("sni", ""),
            "x509_certs": options.get("x509_certs", []),
        },
        "children": [],
    }


dsl = module(
	"dsl",
    compose=_compose,
    discard=_discard,
    dns_lookup_getaddrinfo=_dns_lookup_getaddrinfo,
    domain_name=_domain_name,
    http_connection_tls=_http_connection_tls,
    http_transaction=_http_transaction,
    make_endpoints_for_port=_make_endpoints_for_port,
    new_endpoint=_new_endpoint,
    new_endpoint_pipeline=_new_endpoint_pipeline,
    run=_run,
    tcp_connect=_tcp_connect,
    tls_handshake=_tls_handshake,
)
