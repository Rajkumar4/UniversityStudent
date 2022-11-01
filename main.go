package main

import (
	server "github.com/UnivertsityStudent/Backend"
	"github.com/sirupsen/logrus"
)

 var log = logrus.New()
func main() {
	log.SetLevel(logrus.DebugLevel)
	err := server.InitServer()
	if err!=nil{
		log.Errorf("failed to initalize server %s",err.Error())
		return 
	}
}
