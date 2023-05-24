package runner

import (
	"context"
	"errors"
	"net/http"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/model"
	enginemodel "github.com/ooni/probe-engine/pkg/model"
)

// newSession creates a [model.ExperimentSession] instance.
func (s *State) newSession(
	logger enginemodel.Logger,
	testHelpers map[string][]enginemodel.OOAPIService,
) enginemodel.ExperimentSession {
	return &runnerSession{
		location:    s.location,
		logger:      logger,
		testHelpers: testHelpers,
	}
}

// runnerSession is the [model.ExperimentSession] returned by [State.newSession]
type runnerSession struct {
	location    *model.ProbeLocation
	logger      enginemodel.Logger
	testHelpers map[string][]enginemodel.OOAPIService
}

var _ enginemodel.ExperimentSession = &runnerSession{}

// DefaultHTTPClient implements model.ExperimentSession
func (rs *runnerSession) DefaultHTTPClient() enginemodel.HTTPClient {
	// TODO(bassosimone): stub
	return http.DefaultClient
}

// FetchPsiphonConfig implements model.ExperimentSession
func (rs *runnerSession) FetchPsiphonConfig(ctx context.Context) ([]byte, error) {
	// TODO(bassosimone): stub
	return nil, errors.New("not implemented")
}

// FetchTorTargets implements model.ExperimentSession
func (rs *runnerSession) FetchTorTargets(ctx context.Context, cc string) (map[string]enginemodel.OOAPITorTarget, error) {
	// TODO(bassosimone): stub
	return nil, errors.New("not implemented")
}

// GetTestHelpersByName implements model.ExperimentSession
func (rs *runnerSession) GetTestHelpersByName(name string) ([]enginemodel.OOAPIService, bool) {
	svc, good := rs.testHelpers[name]
	return svc, good
}

// Logger implements model.ExperimentSession
func (rs *runnerSession) Logger() enginemodel.Logger {
	return rs.logger
}

// ProbeCC implements model.ExperimentSession
func (rs *runnerSession) ProbeCC() string {
	// TODO(bassosimone): stub
	return rs.location.IPv4.ProbeCC
}

// ResolverIP implements model.ExperimentSession
func (rs *runnerSession) ResolverIP() string {
	// TODO(bassosimone): stub
	return rs.location.IPv4.ResolverIP
}

// TempDir implements model.ExperimentSession
func (rs *runnerSession) TempDir() string {
	// TODO(bassosimone): stub
	return "/tmp"
}

// TorArgs implements model.ExperimentSession
func (rs *runnerSession) TorArgs() []string {
	// TODO(bassosimone): stub
	return nil
}

// TorBinary implements model.ExperimentSession
func (rs *runnerSession) TorBinary() string {
	// TODO(bassosimone): stub
	return "tor"
}

// TunnelDir implements model.ExperimentSession
func (rs *runnerSession) TunnelDir() string {
	// TODO(bassosimone): stub
	return "/tmp"
}

// UserAgent implements model.ExperimentSession
func (rs *runnerSession) UserAgent() string {
	// TODO(bassosimone): stub
	return "miniooni/0.1.0-dev"
}
