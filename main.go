package main

import (
	"vayeate/database"
	"vayeate/logger"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	log := logger.New("main")
	db := database.NewClient()
	defer db.Close()

	log.Infoln("hello world")
}
