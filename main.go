package main

import (
	"chat/controller"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("unable to load configuration")
	}
	logrus.Info("successfully load configuration")
	ctl := controller.NewController()
	_ = ctl.Run()
}
