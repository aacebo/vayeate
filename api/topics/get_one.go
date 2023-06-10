package topics

import (
	"net/http"
	"vayeate/node"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func GetOne(n *node.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "topic")

		if !n.Topics.Has(name) {
			render.Status(r, 404)
			render.PlainText(w, r, "not found")
			return
		}

		topic := n.Topics.Get(name)
		render.JSON(w, r, map[string]interface{}{
			"name":        topic.Name,
			"created_at":  topic.CreatedAt,
			"subscribers": topic.Subscribers.Len(),
			"messages":    topic.Messages.Len(),
		})
	}
}
