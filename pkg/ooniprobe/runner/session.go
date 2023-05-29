package runner

import (
	"context"
	"errors"
	"net/http"

	"github.com/bassosimone/2023-05-sbs-probe-spec/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/model"
)

// newSession creates a [model.ExperimentSession] instance.
func (s *State) newSession(
	logger model.Logger,
	testHelpers map[string][]model.OOAPIService,
) model.ExperimentSession {
	return &runnerSession{
		location:    s.location,
		logger:      logger,
		testHelpers: testHelpers,
	}
}

// runnerSession is the [model.ExperimentSession] returned by [State.newSession]
type runnerSession struct {
	location    *modelx.ProbeLocation
	logger      model.Logger
	testHelpers map[string][]model.OOAPIService
}

var _ model.ExperimentSession = &runnerSession{}

// DefaultHTTPClient implements model.ExperimentSession
func (rs *runnerSession) DefaultHTTPClient() model.HTTPClient {
	// TODO(bassosimone): stub
	return http.DefaultClient
}

// FetchPsiphonConfig implements model.ExperimentSession
func (rs *runnerSession) FetchPsiphonConfig(ctx context.Context) ([]byte, error) {
	// TODO(bassosimone): stub
	return nil, errors.New("not implemented")
}

// FetchTorTargets implements model.ExperimentSession
func (rs *runnerSession) FetchTorTargets(ctx context.Context, cc string) (map[string]model.OOAPITorTarget, error) {
	// TODO(bassosimone): stub
	return nil, errors.New("not implemented")
}

// GetTestHelpersByName implements model.ExperimentSession
func (rs *runnerSession) GetTestHelpersByName(name string) ([]model.OOAPIService, bool) {
	svc, good := rs.testHelpers[name]
	return svc, good
}

// Logger implements model.ExperimentSession
func (rs *runnerSession) Logger() model.Logger {
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
