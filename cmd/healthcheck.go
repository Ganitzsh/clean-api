package cmd

import (
	"net/http"
	"os"

	"github.com/ganitzsh/f3-te/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var healthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Runs a healthcheck on itself",
	Run: func(cmd *cobra.Command, args []string) {
		api.InitConfig()
		r, err := http.Get(api.Config().GetHostURL() + "/ping")
		if err != nil {
			logrus.Fatalf("Could not query: %v", err)
		}
		if r.StatusCode != http.StatusNoContent {
			logrus.Fatalf("Received invalid status: %d", r.StatusCode)
		}
		os.Exit(0)
	},
}
