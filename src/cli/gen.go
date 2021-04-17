package cli

import "github.com/spf13/cobra"

var cmdGen = &cobra.Command{
	Use: "generate",
	Aliases: []string{ "gen", "g" },
}

func init() {
	cmdroot.AddCommand((cmdGen))
}