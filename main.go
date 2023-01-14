package main

import (
	"vayeate/broker"
	"vayeate/logger"
	"vayeate/utils"
)

func main() {
	log := logger.New("main")
	port, err := utils.GetPort()

	if err != nil {
		log.Error(err)
	}

	b := broker.New(port)
	log.Infof("listening on port %d", port)
	b.Listen()

	defer b.Close()
}
