package modelx

import "fmt"

// ASNumber is an autonomous system number.
type ASNumber int64

// String converts the AS number to the "AS%d" string.
func (asn ASNumber) String() string {
	return fmt.Sprintf("AS%d", asn)
}
