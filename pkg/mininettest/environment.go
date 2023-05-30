package mininettest

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/dslx"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
)

// Environment is the environment in which we run mininettest. The zero value
// of this structure is invalid; use [NewEnvironment] to instantiate.
type Environment struct {
	// idGenerator is used to generate identifiers.
	idGenerator *atomic.Int64

	// logger is the logger to use.
	logger enginemodel.Logger

	// zeroTime is the reference time of the measurement.
	zeroTime time.Time
}

// NewEnvironment creates a new [Environment].
func NewEnvironment(logger enginemodel.Logger, zeroTime time.Time) *Environment {
	return &Environment{
		idGenerator: &atomic.Int64{},
		logger:      logger,
		zeroTime:    zeroTime,
	}
}

// ErrNoSuchNettestlet indicates a given mininettest does not exist.
var ErrNoSuchNettestlet = errors.New("mininettest: no such mininettest")

// Run runs the given mininettest in the current goroutine. This function only
// returns an error only in case a fundamental error has occurred (e.g., not
// being able to parse the descriptor With field).
func (env *Environment) Run(
	ctx context.Context,
	descr *modelx.MiniNettestDescriptor,
) (*dslx.Observations, error) {
	switch descr.Run {
	case "dns-lookup@v1":
		return env.dnsLookupV1Main(ctx, descr)

	case "http-address@v1":
		return env.httpAddressV1Main(ctx, descr)

	case "http-domain@v1":
		return env.httpDomainV1Main(ctx, descr)

	case "https-domain@v1":
		return env.httpsDomainV1Main(ctx, descr)

	case "tcp-connect-address@v1":
		return env.tcpConnectAddressV1Main(ctx, descr)

	case "tcp-connect-domain@v1":
		return env.tcpConnectDomainV1Main(ctx, descr)

	default:
		return nil, ErrNoSuchNettestlet
	}
}
