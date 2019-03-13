package main

import (
	"github.com/ganitzsh/f3-te/api"
	"github.com/ganitzsh/f3-te/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	cmd.Execute(func() {
		api.InitConfig()
		api.InitStore()
		if err := api.Start(); err != nil {
			logrus.Fatalf("Could not run the server: %v", err)
		}
	})
}
