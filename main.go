package main

import (
	"github.com/ganitzsh/f3-te/api"
	"github.com/sirupsen/logrus"
)

func init() {
	api.InitConfig()
	// api.InitStore(api.NewDocumentInMemStore())
}

func main() {
	if err := api.Start(); err != nil {
		logrus.Fatalf("Could not run the server: %v", err)
	}
}
