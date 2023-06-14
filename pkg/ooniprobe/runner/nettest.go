package runner

//
// nettest.go contains code to create nettests.
//
// The nettest is the user facing executable network experiment
// interface, while experiment is the corresponding implementation
// inside of the OONI probe engine. We will eventually refactor
// the probe engine to merge nettests and experiments.
//

import (
	"context"
	"errors"

	"github.com/ooni/2023-05-richer-input/pkg/modelx"
)

// errNoSuchNettest indicates that the given nettest does not exist
var errNoSuchNettest = errors.New("no such nettest")

// nettestFactory constructs a nettest.
type nettestFactory = func(args *modelx.InterpreterNettestRunArguments,
	config *modelx.InterpreterConfig, ix *Interpreter) (nt nettest, err error)

// nettestRegistry maps nettests to their constructors.
var nettestRegistry = map[string]nettestFactory{
	"facebook_messenger": fbmessengerNew,
	"signal":             signalNew,
	"telegram":           telegramNew,
	"urlgetter":          urlgetterNew,
	"web_connectivity":   webconnectivityNew,
	"whatsapp":           whatsappNew,
}

// newNettest creates a new [nettest] instance.
func newNettest(args *modelx.InterpreterNettestRunArguments,
	config *modelx.InterpreterConfig, ix *Interpreter) (nettest, error) {
	factory := nettestRegistry[args.NettestName]
	if factory == nil {
		return nil, errNoSuchNettest
	}
	return factory(args, config, ix)
}

// nettest is a nettest instance.
type nettest interface {
	Run(ctx context.Context) error
}
