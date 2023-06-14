package runner

//
// session.go contains code to create a model.ExperimentSession. This is a data
// type required by the current OONI probe engine to execute experiments.
//

import (
	"context"
	"errors"
	"net/http"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
	"github.com/ooni/probe-engine/pkg/model"
)

// newSession creates a new [model.ExperimentSession] instance.
func newSession(
	location modelx.InterpreterLocation,
	logger model.Logger,
	testHelpers map[string][]model.OOAPIService,
) (model.ExperimentSession, error) {
	v4Location := location.IPv4()
	if v4Location.IsNone() {
		return nil, ErrMissingIPv4Location
	}

	sess := &session{
		logger:      logger,
		testHelpers: testHelpers,
		v4Location:  v4Location.Unwrap(),
	}

	return sess, nil
}

// session is the [model.ExperimentSession] returned by [State.newSession]
type session struct {
	logger      model.Logger
	testHelpers map[string][]model.OOAPIService
	v4Location  *modelx.Location
}

var _ model.ExperimentSession = &session{}

// DefaultHTTPClient implements model.ExperimentSession
func (s *session) DefaultHTTPClient() model.HTTPClient {
	// TODO(bassosimone): stub
	return http.DefaultClient
}

// FetchPsiphonConfig implements model.ExperimentSession
func (s *session) FetchPsiphonConfig(ctx context.Context) ([]byte, error) {
	// TODO(bassosimone): stub
	return nil, errors.New("not implemented")
}

// FetchTorTargets implements model.ExperimentSession
func (s *session) FetchTorTargets(ctx context.Context, cc string) (map[string]model.OOAPITorTarget, error) {
	// TODO(bassosimone): stub
	return nil, errors.New("not implemented")
}

// GetTestHelpesByName implements model.ExperimentSession
func (s *session) GetTestHelpersByName(name string) ([]model.OOAPIService, bool) {
	svc, good := s.testHelpers[name]
	return svc, good
}

// Logger implements model.ExperimentSession
func (s *session) Logger() model.Logger {
	return s.logger
}

// ProbeCC implements model.ExperimentSession
func (s *session) ProbeCC() string {
	return s.v4Location.ProbeCC
}

// ResolverIP implements model.ExperimentSession
func (s *session) ResolverIP() string {
	return s.v4Location.ResolverIP
}

// TempDir implements model.ExperimentSession
func (s *session) TempDir() string {
	// TODO(bassosimone): stub
	return "/tmp"
}

// TorArgs implements model.ExperimentSession
func (s *session) TorArgs() []string {
	// TODO(bassosimone): stub
	return nil
}

// TorBinary implements model.ExperimentSession
func (s *session) TorBinary() string {
	// TODO(bassosimone): stub
	return "tor"
}

// TunnelDir implements model.ExperimentSession
func (s *session) TunnelDir() string {
	// TODO(bassosimone): stub
	return "/tmp"
}

// UserAgent implements model.ExperimentSession
func (s *session) UserAgent() string {
	// TODO(bassosimone): stub
	return "miniooni/0.1.0-dev"
}
