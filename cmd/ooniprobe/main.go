package main

import (
	"github.com/ooni/probe-engine/pkg/runtimex"
	"github.com/spf13/cobra"
)

func main() {
	// create the root command
	root := &cobra.Command{
		Use: "ooniprobe",
	}

	// create the runx command
	root.AddCommand(newRunxSubcommand())

	// execute the selected subcommand
	runtimex.Try0(root.Execute())
}
