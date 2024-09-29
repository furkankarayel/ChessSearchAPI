package scoutfish

import (
	"engine"
	"net/http"
)

type ScoutfishHandler struct {
	scoutfish *Scoutfish
}

func New() *engine.Route {
	scoutfishInstance := NewScoutfish("LumbrasGigaBase-2020")
	scoutfishHandler := &ScoutfishHandler{scoutfish: scoutfishInstance}
	return &engine.Route{
		WithLogger: true,
		Handler:    scoutfishHandler,
	}
}

func (s *ScoutfishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	switch head {
	case "":
		s.home(w, r)
	case "test":
		s.test(w, r)
	default:
		engine.Respond(w, r, http.StatusNotFound, "scoutfish path not found")
	}
}
