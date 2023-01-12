package main

import (
	"vayeate/broker"
	"vayeate/database"
	"vayeate/frame"
	"vayeate/logger"
	"vayeate/utils"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	log := logger.New("main")

	db := database.New()
	defer db.Close()

	port, err := utils.GetPort()

	if err != nil {
		log.Error(err)
	}

	b := broker.New(port)
	log.Infof("listening on port %d", port)
	b.Listen(func(f *frame.Frame) {
		log.Infof("new frame with length %d", len(f.Payload))
	})

	defer b.Close()
}
