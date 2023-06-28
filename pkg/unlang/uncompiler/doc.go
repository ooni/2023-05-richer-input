/*
Package uncompiler contains the micro-nettest-language compiler. Typically one loads
a [*ASTNode] encoded as JSON using [encoding/json]. Then, one creates a [*Compiler] using
[NewCompiler] and calls [*Compiler.Compile] to obtain a [unruntime.Func], and executes
the given [unruntime.Func]. One can use [*Compiler.RegisterFuncTemplate] to register new
[unruntime.Func] beyond the built-in ones defined by the [unruntime] package and automatically
loaded by the [NewCompiler] factory. This package defines a [FuncTemplate] for each
built-in [unruntime] function loaded by [NewCompiler].
*/
package uncompiler
