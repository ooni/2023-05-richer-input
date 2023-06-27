package ridsl

import "fmt"

// Dump prints debugging information on the standard output.
func Dump(f *Func) {
	doDump(f, 0)
}

// doDump is the worker function called by [Dump].
func doDump(f *Func, indent int) {
	fmt.Printf(
		"%s%s :: %s -> %s\n",
		dumpIndent(indent),
		f.Name,
		f.InputType,
		f.OutputType,
	)
	for _, child := range f.Children {
		doDump(child, indent+1)
	}
}

// dumpIndent converts the indent value from integer to a string.
func dumpIndent(value int) (out string) {
	for ; value > 0; value-- {
		out += "    "
	}
	return
}
