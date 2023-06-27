package riengine

// DNSLookupInputArguments contains the arguments for "dns_lookup_input".
type DNSLookupInputArguments struct {
	Domain string `json:"domain"`
}

// DNSLookupStaticArguments contains the arguments for "dns_lookup_static".
type DNSLookupStaticArguments struct {
	Addresses []string `json:"addresses"`
}

// DNSLookupUDPArguments contains the arguments for "dns_lookup_udp".
type DNSLookupUDPArguments struct {
	Endpoint string `json:"endpoint"`
}
