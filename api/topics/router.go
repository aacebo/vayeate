package topics

import (
	"vayeate/node"

	"github.com/go-chi/chi/v5"
)

func NewRouter(r chi.Router, n *node.Node) {
	r.Get("/topics", Get(n))
	r.Get("/topics/{topic}", GetOne(n))
}
