package main

import (
	"github.com/sirupsen/logrus"
)

func init() {
	initConfig()
}

func main() {
	api := NewPaymentAPI()
	if err := api.Start(); err != nil {
		logrus.Fatal(err)
	}
}
