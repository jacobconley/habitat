package cli

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{

	Use:   "habitat",
	Short: "The root command for the Habitat CLI.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CLI test")
	},
}

func RunCLI() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}