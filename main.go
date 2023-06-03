package main

import (
	"vayeate/common"
	"vayeate/logger"
	"vayeate/node"
)

func main() {
	log := logger.New("vayeate")
	n, err := node.New(
		common.GetEnv("VAYEATE_CLIENT_PORT", "6789"),
		common.GetEnv("VAYEATE_USERNAME", "admin"),
		common.GetEnv("VAYEATE_PASSWORD", "admin"),
	)

	if err != nil {
		log.Error(err)
		return
	}

	defer n.Close()
	err = n.Listen()

	if err != nil {
		log.Error(err)
	}
}
