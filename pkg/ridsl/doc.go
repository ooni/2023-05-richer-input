// Package ridsl contains the richer-input DSL.
//
// # Typical Usage
//
// When using [ridsl], your objective is to use [Func] constructors such as [DNSLookupGetaddrinfo]
// and [TCPConnect] and compose them together using [Compose] and other functions to create a
// [Func] representing a nettest. When you have created such a [Func] you can [Dump] it to inspect
// its types and [Compile] it to obtain an [ASTNode]. The [ASTNode] serialies to [encoding/json]
// using a data format that is compatible with the one expected by the [riengine].
//
// # Writing Extensions
//
// An extension is a Go func that returns a properly configured [Func]. Typically, you need to
// wrote extensions to write filters that intercept the results of calling another [Func] and
// use this information to write specific test keys for a nettest.
//
// # Type differences with riengine
//
// We document the input and output types of each [Func]. The runtime implementation inside
// [riengine], however, uses extended types. In particular, a type T defined by this package
// becomes the sum type T|error|*Exception|*Skip in the [riengine] package. Because we extend
// each type equally, we have chosen to omit the additional error|*Exception|*Skip in the
// documentation and we do not use it for type checking. However, the [Dump] function shows
// the full [riengine] type including the error|*Exception|*Skip alternatives. Likewise, the
// [CompleteTypeName] function maps a type T to its full [riengine] type.
package ridsl
