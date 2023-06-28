/*
Package unlang supports the micro-nettest language. This language allows one to write micro
nettests. A nettest is a network test run by OONI Probe. A micro nettest is a set of instructions
for performing a nettest. For example, the Facebook Messenger OONI nettest measures several
services crucial to Facebook Messenger operations as part of a single nettest run. One can express
the measurement of each Facebook Messenger service using the micro-nettest language.

# Domain-Specific Language

The undsl subpackage contains a domain-specific language for writing micro nettests. This
language allows one to express measuring specific targets in terms of predefined network
operations, such as DNS lookup, TCP connect, TLS handshake, and QUIC handshake. Once one has
defined how to perform the measurement, one can export the abstract syntax tree (AST) describing the
measurement. The expected usage pattern is that one serializes such an AST to JSON and,
in turn, serves the JSON to the OONI probes, which will interpret it and execute the measurement.

# Compiler

The uncompiler subpackage compiles the JSON-serialized AST to runtime data structures.

# Runtime

The unruntime subpackage implements the micro-nettest-language runtime. One can use such a
runtime directly or indirectly. Typically, one compiles an AST to runtime structures using the
compiler and then executes the measurement; however, it is also possible to manually
compose the basic runtime primitives, thus bypassing the dsl and the compiler.
*/
package unlang
