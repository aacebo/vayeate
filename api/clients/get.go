package clients

import (
	"net/http"
	"vayeate/node"

	"github.com/go-chi/render"
)

func Get(n *node.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, n.Clients.Slice())
	}
}
