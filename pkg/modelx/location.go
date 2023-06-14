package modelx

import "fmt"

// LocationASNumber is an autonomous system number.
type LocationASNumber int64

// String converts the AS number to the "AS%d" string.
func (asn LocationASNumber) String() string {
	return fmt.Sprintf("AS%d", asn)
}

// Location is the location relative to a given IP address.
type Location struct {
	// ProbeIP is the probe IP address.
	ProbeIP string `json:"probe_ip"`

	// ProbeASN is the probe IP address AS number.
	ProbeASN LocationASNumber `json:"probe_asn"`

	// ProbeCC is the probe IP country code.
	ProbeCC string `json:"probe_cc"`

	// ProbeNetworkName is the probe IP network name.
	ProbeNetworkName string `json:"probe_network_name"`

	// ResolverIP is the IP address used by getaddrinfo.
	ResolverIP string `json:"resolver_ip"`

	// ResolverASN is the resolver IP AS number.
	ResolverASN LocationASNumber `json:"resolver_asn"`

	// ResolverCC is the resolver IP country code.
	ResolverCC string `json:"resolver_cc"`

	// ResolverNetworkName is the resolver IP network name.
	ResolverNetworkName string `json:"resolver_network_name"`
}

// SameASNAndCC returns whether two locations have the same ASN and CC.
func (a *Location) SameASNAndCC(b *Location) bool {
	return a.ProbeASN == b.ProbeASN && a.ProbeCC == b.ProbeCC
}
