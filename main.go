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

	b, err := broker.New(port)

	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("listening on port %d", port)
	err = b.Listen()

	if err != nil {
		log.Error(err)
		return
	}

	defer b.Close()
}
