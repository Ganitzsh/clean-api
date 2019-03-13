package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Simple Payment API",
	Long:  "",
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(healthcheckCmd)
	rootCmd.AddCommand(docgenCmd)
}

func Execute(mainFunc func()) {
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		mainFunc()
	}
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
