package scoutfish

import (
	"engine"
	"net/http"
)

type ScoutfishHandler struct{}

func New() *engine.Route {
	scoutfish := &ScoutfishHandler{}
	return &engine.Route{
		WithLogger: true,
		Handler:    scoutfish,
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
