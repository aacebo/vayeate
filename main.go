package main

import (
	"vayeate/env"
	"vayeate/frame"
	"vayeate/logger"
	"vayeate/queue"
	"vayeate/server"
)

var log = logger.New("main")

func main() {
	port, err := env.GetPort()

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
		log.Infof("new socket <%s>", s.GetRemoteAddress())
	}

	var onFrame = func(s *server.Socket, f *frame.Frame) {
		if f.IsAssert() {
			onAssert(serv, s, f)
		} else if f.IsProduce() {
			onProduce(serv, s, f)
		} else if f.IsConsume() {
			onConsume(serv, s, f)
		}
	}

	err = serv.Listen(onConnect, onFrame, nil)

	if err != nil {
		log.Error(err)
		return
	}

	defer serv.Close()
}

func onAssert(serv *server.Server, s *server.Socket, f *frame.Frame) {
	q := serv.AddQueue(queue.New(f.GetSubject()))
	log.Infof("assert %s", q.Name)
}

func onProduce(serv *server.Server, s *server.Socket, f *frame.Frame) {
	qs := serv.GetQueues(f.GetSubject())

	for _, q := range qs {
		q.Push(f.Body)
		log.Infof("produce %s => %s", q.Name, f.GetBody())
	}
}

func onConsume(serv *server.Server, s *server.Socket, f *frame.Frame) {
	qs := serv.GetQueues(f.GetSubject())

	for _, q := range qs {
		q.Consume(s.ID)
		log.Infof("consume %s => %s", s.ID, q.Name)
	}
}
