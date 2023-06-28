// Package ridsl contains the richer-input DSL.
//
// # Typical Usage
//
// When using [ridsl], your objective is to use [*Func] constructors such as [DNSLookupGetaddrinfo]
// and [TCPConnect] and compose them together using [Compose] and other functions to create a
// [*Func] representing a nettest. When you have created such a [*Func] you can [Dump] it to inspect
// its types and [Compile] it to obtain an [*ASTNode]. The [*ASTNode] serializes to [encoding/json]
// using a data format that is compatible with the one expected by the [riengine].
//
// # Writing Extensions
//
// An extension is a Go func that returns a properly configured [*Func]. Typically, you need to
// write extensions to write filters that intercept the results of calling another [*Func] and
// use this information to write specific test keys for a nettest.
//
// # Main and Complete Type
//
// All the [riengine] functions handle the [SumType] of [ErrorType], [ExceptionType], and
// [SkipType]. Each of these three types represents a well defined abnormal condition that needs to
// be handled separately. For brevity, we are not going to document these three types for
// each function defined by this package. Rather, we will only document the "main" input and
// output type of each function, defined as the actual type minus the [SumType] of [ErrorType],
// [ExceptionType], and [SkipType]. Conversely, we define "complete type" the main type plus
// the [SumType] of [ErrorType], [ExceptionType], and [SkipType].
//
// For example, the [TCPConnect] [Func] has main type:
//
//	EndpointType -> TCPConnectionType
//
// and its complete type is:
//
//	EndpointType|ErrorType|ExceptionType|SkipType -> TCPConnectionType|ErrorType|ExceptionType|SkipType
//
// The [CompleteType] function takes in input the main type and returns the complete type. The
// [MainType] function takes in input the complete type and returns the main type.
package ridsl
