package main

import (
	"vayeate/frame"
	"vayeate/logger"
	"vayeate/queue"
	"vayeate/server"
	"vayeate/utils"
)

var log = logger.New("main")

func main() {
	port, err := utils.GetPort()

	if err != nil {
		log.Error(err)
	}

	serv, err := server.New(port)

	if err != nil {
		log.Error(err)
		return
	}

	log.Infof("listening on port %d", port)

	var onConnect = func(s *server.Socket) {
		log.Info(s)
	}

	var onFrame = func(s *server.Socket, f *frame.Frame) {
		if f.IsAssert() {
			onAssert(serv, s, f)
		} else if f.IsProduce() {
			onProduce(serv, s, f)
		}
	}

	err = serv.Listen(onConnect, onFrame)

	if err != nil {
		log.Error(err)
		return
	}

	defer serv.Close()
}

func onAssert(server *server.Server, s *server.Socket, f *frame.Frame) {
	q := server.AddQueue(queue.New(f.GetSubject()))
	log.Info(q.ID)
}

func onProduce(server *server.Server, s *server.Socket, f *frame.Frame) {
	qs := server.GetQueues(f.GetSubject())

	for _, q := range qs {
		q.Push(f.Body)
	}
}
