package undsl

import "github.com/ooni/2023-05-richer-input/pkg/unlang/uncompiler"

// HTTPTransactionOption is an option for the [HTTPTransaction].
type HTTPTransactionOption func(c *uncompiler.HTTPTransactionArguments)

// HTTPTransactionOptionAccept sets the Accept header.
func HTTPTransactionOptionAccept(value string) HTTPTransactionOption {
	return func(c *uncompiler.HTTPTransactionArguments) {
		c.AcceptHeader = value
	}
}

// HTTPTransactionOptionAcceptLanguage sets the Accept-Language header.
func HTTPTransactionOptionAcceptLanguage(value string) HTTPTransactionOption {
	return func(c *uncompiler.HTTPTransactionArguments) {
		c.AcceptLanguageHeader = value
	}
}

// HTTPTransactionOptionHost sets the Host header.
func HTTPTransactionOptionHost(value string) HTTPTransactionOption {
	return func(c *uncompiler.HTTPTransactionArguments) {
		c.HostHeader = value
	}
}

// HTTPTransactionOptionMethod sets the method.
func HTTPTransactionOptionMethod(value string) HTTPTransactionOption {
	return func(c *uncompiler.HTTPTransactionArguments) {
		c.RequestMethod = value
	}
}

// HTTPTransactionOptionResponseBodySnapshotSize sets the maximum response body snapshot size.
func HTTPTransactionOptionResponseBodySnapshotSize(value int) HTTPTransactionOption {
	return func(c *uncompiler.HTTPTransactionArguments) {
		c.ResponseBodySnapshotSize = value
	}
}

// HTTPTransactionOptionReferer sets the referer.
func HTTPTransactionOptionReferer(value string) HTTPTransactionOption {
	return func(c *uncompiler.HTTPTransactionArguments) {
		c.RefererHeader = value
	}
}

// HTTPTransactionOptionURLHost sets the URL host.
func HTTPTransactionOptionURLHost(value string) HTTPTransactionOption {
	return func(c *uncompiler.HTTPTransactionArguments) {
		c.URLHost = value
	}
}

// HTTPTransactionOptionURLPath sets the URL path.
func HTTPTransactionOptionURLPath(value string) HTTPTransactionOption {
	return func(c *uncompiler.HTTPTransactionArguments) {
		c.URLPath = value
	}
}

// HTTPTransactionOptionURLScheme sets the URL scheme.
func HTTPTransactionOptionURLScheme(value string) HTTPTransactionOption {
	return func(c *uncompiler.HTTPTransactionArguments) {
		c.URLScheme = value
	}
}

// HTTPTransactionOptionUserAgent sets the User-Agent header.
func HTTPTransactionOptionUserAgent(value string) HTTPTransactionOption {
	return func(c *uncompiler.HTTPTransactionArguments) {
		c.UserAgentHeader = value
	}
}

// HTTPTransaction returns a [*Func] that sends a request and receives a response.
//
// The main returned [*Func] type is: [TCPConnectionType] | [TLSConnectionType] | [QUICConnectionType] -> [VoidType].
func HTTPTransaction(options ...HTTPTransactionOption) *Func {
	args := &uncompiler.HTTPTransactionArguments{}
	for _, option := range options {
		option(args)
	}
	return &Func{
		Name:       templateName[uncompiler.HTTPTransactionTemplate](),
		InputType:  SumType(TCPConnectionType, TLSConnectionType, QUICConnectionType),
		OutputType: VoidType,
		Arguments:  args,
		Children:   []*Func{},
	}
}
