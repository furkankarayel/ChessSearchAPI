package pgnextract

import (
	"engine"
	"net/http"
)

type PgnextracthHandler struct{}

func New() *engine.Route {
	pgnextract := &PgnextracthHandler{}
	return &engine.Route{
		WithLogger: true,
		Handler:    pgnextract,
	}
}

func (p *PgnextracthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	switch head {
	case "":
		p.home(w, r)
	case "test":
		p.test(w, r)
	default:
		engine.Respond(w, r, http.StatusNotFound, "pgn-extract path not found")
	}
}
