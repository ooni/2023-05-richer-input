package model

// ProbeLocation is the location of a probe.
type ProbeLocation struct {
	// IPv4 contains the IPv4 location.
	IPv4 *ProbeLocationIPAddr `json:"ipv4"`

	// IPv6 contains the IPv6 location.
	IPv6 *ProbeLocationIPAddr `json:"ipv6"`
}

// ProbeLocationIPAddr is the location relative to a given IP address.
type ProbeLocationIPAddr struct {
	// ProbeIP is the probe IP address.
	ProbeIP string `json:"probe_ip"`

	// ProbeASN is the probe IP address AS number.
	ProbeASN ASNumber `json:"probe_asn"`

	// ProbeCC is the probe IP country code.
	ProbeCC string `json:"probe_cc"`

	// ProbeNetworkName is the probe IP network name.
	ProbeNetworkName string `json:"probe_network_name"`

	// ResolverIP is the IP address used by getaddrinfo.
	ResolverIP string `json:"resolver_ip"`

	// ResolverASN is the resolver IP AS number.
	ResolverASN ASNumber `json:"resolver_asn"`

	// ResolverCC is the resolver IP country code.
	ResolverCC string `json:"resolver_cc"`

	// ResolverNetworkName is the resolver IP network name.
	ResolverNetworkName string `json:"resolver_network_name"`
}
