package topics

import (
	"fmt"
	"net/http"
	"vayeate/node"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func GetOne(n *node.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "topic_name")
		fmt.Println(name)

		if !n.Topics.Has(name) {
			render.Status(r, 404)
			render.PlainText(w, r, "not found")
			return
		}

		topic := n.Topics.Get(name)
		render.JSON(w, r, topic)
	}
}
