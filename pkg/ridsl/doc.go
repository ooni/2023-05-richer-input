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
package ridsl
