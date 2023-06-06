package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"vayeate/api"
	"vayeate/common"
	"vayeate/logger"
	"vayeate/node"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {
	var wg sync.WaitGroup
	log := logger.New("vayeate")

	wg.Add(2)

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

	go func() {
		err = n.Listen()

		if err != nil {
			log.Error(err)
		}

		wg.Done()
	}()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.AllowAll().Handler)
	r.Mount("/", api.NewRouter(n))

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", common.GetEnv("VAYEATE_API_PORT", "8080")), r)

		if err != nil {
			log.Error(err)
		}

		wg.Done()
	}()

	wg.Wait()
}
