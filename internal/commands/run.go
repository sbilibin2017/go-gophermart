package commands

import "github.com/spf13/cobra"

func Run(cmd *cobra.Command) int {
	err := cmd.Execute()
	if err != nil {
		return 1
	}
	return 0
}
