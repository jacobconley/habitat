package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var cmdroot = &cobra.Command{

	Use:   "habitat",
	Short: "The root command for the Habitat CLI.",

	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("CLI test")
	// },
}

func RunCLI() {
	if err := cmdroot.Execute(); err != nil {
		os.Exit(1)
	}
}