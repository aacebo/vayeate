package api

import (
	"net/http"
	"vayeate/api/clients"
	"vayeate/node"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func NewRouter(n *node.Node) *chi.Mux {
	r := chi.NewRouter()
	clients.NewRouter(r, n)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, n)
	})

	return r
}
