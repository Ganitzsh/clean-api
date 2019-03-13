package cmd

import (
	"fmt"

	"github.com/ganitzsh/f3-te/api"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCode name: %s\n", api.Version, api.ReleaseName)
	},
}
