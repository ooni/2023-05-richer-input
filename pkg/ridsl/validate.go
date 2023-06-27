package ridsl

import (
	"fmt"
	"math"
	"net"
	"strconv"

	"github.com/ooni/probe-engine/pkg/runtimex"
)

// panicUnlessValidIPAddress panics unless address is a valid IP address.
func panicUnlessValidIPAddress(address string) {
	runtimex.Assert(net.ParseIP(address) != nil, fmt.Sprintf("%s: not a valid IP address", address))
}

// panicUnlessValidPort panics unless port is a valid port.
func panicUnlessValidPort(sport string) {
	nport := runtimex.Try1(strconv.Atoi(sport))
	runtimex.Assert(nport >= 0 && nport <= math.MaxUint16, fmt.Sprintf("%s: not a valid port", sport))
}
