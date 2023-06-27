package riengine

// TLSHandshakeArguments contains arguments for "tls_handshake".
type TLSHandshakeArguments struct {
	ALPN       []string `json:"alpn,omitempty"`
	SkipVerify bool     `json:"skip_verify,omitempty"`
	SNI        string   `json:"sni,omitempty"`
	X509Certs  []string `json:"x509_certs,omitempty"`
}
