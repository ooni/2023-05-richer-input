package uncompiler_test

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"
	"github.com/ooni/2023-05-richer-input/pkg/unlang/unruntime"
	"github.com/ooni/probe-engine/pkg/runtimex"
)

func ExampleCompiler_Compile() {
	// create the compiler we're going to use
	comp := uncompiler.NewCompiler()

	// create a fake AST to compile (usually we'll unmarshal this from JSON)
	astRoot := &uncompiler.ASTNode{
		Func:      "compose",
		Arguments: []byte(`{}`),
		Children: []*uncompiler.ASTNode{{
			Func:      "tcp_connect",
			Arguments: []byte(`{}`),
			Children:  []*uncompiler.ASTNode{},
		}, {
			Func:      "tls_handshake",
			Arguments: []byte(`{}`),
			Children:  []*uncompiler.ASTNode{},
		}},
	}

	// compile the AST
	f0, err := comp.Compile(astRoot)

	// ignore the function
	_ = f0

	// make sure there was no error
	fmt.Printf("%+v\n", err)
	// output: <nil>
}

// CounterFunc is the [unruntime.Func] we're adding in the middle.
type CounterFunc struct {
	// ok counts the number of successes
	ok *atomic.Int64

	// fail counts the number of failures
	fail *atomic.Int64
}

// Apply implements unruntime.Func.
func (f *CounterFunc) Apply(ctx context.Context, rtx *unruntime.Runtime, input any) (output any) {
	// type switch over the input type
	switch xinput := input.(type) {

	// handle the cases in which we should return immediately
	case *unruntime.Skip:
		return xinput
	case *unruntime.Exception:
		return xinput

	// on success, register the success and continue
	case *unruntime.TCPConnection:
		f.ok.Add(1)
		return xinput

	// on failure, register the failure and stop processing because we don't want
	// a subsequent filter to be tripped by an error *we already handled*
	case error:
		f.fail.Add(1)
		return &unruntime.Skip{}

	default:
		return unruntime.NewException("unexpected type %T", input)
	}
}

// CounterFuncTemplate is the template for [CounterFunc].
type CounterFuncTemplate struct {
	// ok counts the number of successes
	ok *atomic.Int64

	// fail counts the number of failures
	fail *atomic.Int64
}

// Compile implements uncompiler.FuncTemplate.
func (t CounterFuncTemplate) Compile(compiler *uncompiler.Compiler, node *uncompiler.ASTNode) (unruntime.Func, error) {
	return &CounterFunc{t.ok, t.fail}, nil
}

// TemplateName implements uncompiler.FuncTemplate.
func (t CounterFuncTemplate) TemplateName() string {
	return "counter"
}

func ExampleCompiler_RegisterFuncTemplate() {
	// create the compiler we're going to use
	comp := uncompiler.NewCompiler()

	// create the OK and fail counters
	ok := &atomic.Int64{}
	fail := &atomic.Int64{}

	// register the new template
	comp.RegisterFuncTemplate(CounterFuncTemplate{ok, fail})

	// create AST manually (usually it would be decoded from JSON)
	astRoot := &uncompiler.ASTNode{
		Func:      "compose",
		Arguments: []byte(`{}`),
		Children: []*uncompiler.ASTNode{{
			Func:      "tcp_connect",
			Arguments: []byte(`{}`),
			Children:  []*uncompiler.ASTNode{},
		}, {
			Func:      "counter",
			Arguments: []byte{},
			Children:  []*uncompiler.ASTNode{},
		}, {
			Func:      "tls_handshake",
			Arguments: []byte(`{}`),
			Children:  []*uncompiler.ASTNode{},
		}},
	}

	// compile the AST
	f0 := runtimex.Try1(comp.Compile(astRoot))

	// create endpoint for the AST
	epnt := &unruntime.Endpoint{
		Address: "8.8.8.8:443",
		Domain:  "dns.google.com",
	}

	// create a runtime
	rtx := unruntime.NewRuntime()

	// execute the f0 function
	result := f0.Apply(context.Background(), rtx, epnt)

	// make sure the handshake succeded
	switch result.(type) {
	case *unruntime.TLSConnection:
		// what we expect

	default:
		panic(fmt.Errorf("expected to see a TLSConnection here"))
	}

	// check the counters.
	fmt.Printf("%d %d\n", ok.Load(), fail.Load())
	// output: 1 0
}
