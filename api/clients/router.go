package clients

import (
	"vayeate/node"

	"github.com/go-chi/chi/v5"
)

func NewRouter(r chi.Router, n *node.Node) {
	r.Get("/clients", Get(n))
}
