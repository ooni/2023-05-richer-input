/*
Package undsl contains the micro-nettest-language DSL. This DSL allows one to define a
measurement algorithm by composing several [*Func] together. The result of this composition
is a single [*Func], which one can export to [*ASTNode] and serialize to JSON using the
[encoding/json] package, to produce input for the [uncompiler] package.

Composing [*Func] using [Compose] or using other package functions that combine existing
[*Func] together is subject to type checking. If type checking fails, the code will PANIC
with an explanatory error message. The documentation of each [*Func] specifies
the main input and output types it accepts and returns. We define as "main
type" the type accepted or returned by a [unruntime] function under normal operating
conditions. However, all [unruntime] functions accept and return the [SumType]
of the main type and of the following three types:

- error

- [*unruntime.Exception]

- [*unruntime.Skip]

Because each [unruntime] function always accepts and returns the [SumType] of this three types in
and the main type, the documentation does not explicitly mention these three types, for brevity.
However, the [Dump] function shows the complete types accepted and returned by each function. Use
[CompleteType] to convert the main type to the complete type; use [MainType] to convert the complete
type back to its main type.
*/
package undsl
