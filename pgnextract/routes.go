package pgnextract

import (
	"engine"

	"net/http"
)

type PgnextracthHandler struct {
	pgnextract *Pgnextract
}

func New() *engine.Route {
	pgnInstance := NewPgnextract("LumbrasGigaBase-2020")
	pgnextractHandler := &PgnextracthHandler{pgnextract: pgnInstance}
	return &engine.Route{
		WithLogger: true,
		Handler:    pgnextractHandler,
	}
}

func (p *PgnextracthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	switch head {
	case "":

	case "player":
		p.searchPlayer(w, r)
	default:
		engine.Respond(w, r, http.StatusNotFound, "pgn-extract path not found")
	}
}
