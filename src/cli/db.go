package cli

import "github.com/spf13/cobra"

// [ISSUE #7] Maybe we could specify which db at this level here and leave all of its child commands the same

var cmdDB = &cobra.Command{
	Use: "db",
}

func init() { 
	cmdroot.AddCommand(cmdDB)	
}