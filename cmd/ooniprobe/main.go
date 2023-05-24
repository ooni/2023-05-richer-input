package main

import (
	"github.com/ooni/probe-engine/pkg/runtimex"
	"github.com/spf13/cobra"
)

var (
	// verbose controls the verbosity
	verbose bool
)

func main() {
	// create the root command
	root := &cobra.Command{
		Use: "ooniprobe",
	}

	// add the -v, --verbose command line flag
	root.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"enable verbose logging",
	)

	// create the runx command
	root.AddCommand(newRunxSubcommand())

	// execute the selected subcommand
	runtimex.Try0(root.Execute())
}
