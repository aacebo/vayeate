package main

import (
	"vayeate/logger"
	"vayeate/server"
	"vayeate/utils"
)

func main() {
	log := logger.New("main")
	port, err := utils.GetPort()

	if err != nil {
		log.Error(err)
	}

	s, err := server.New(port)

	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("listening on port %d", port)
	err = s.Listen()

	if err != nil {
		log.Error(err)
		return
	}

	defer s.Close()
}
