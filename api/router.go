package api

import (
	"vayeate/api/clients"
	"vayeate/node"

	"github.com/go-chi/chi/v5"
)

func NewRouter(n *node.Node) *chi.Mux {
	r := chi.NewRouter()
	clients.NewRouter(r, n)
	return r
}
