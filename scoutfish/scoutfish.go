package scoutfish

import (
	"engine"
	"engine/helper"
	"log"
	"net/http"
	"strconv"
)

type Scoutfish struct{}

func New() *engine.Route {
	scoutfish := &Scoutfish{}
	return &engine.Route{
		WithLogger: true,
		Handler:    scoutfish,
	}
}

func (s *Scoutfish) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = engine.ShiftPath(r.URL.Path)
	switch head {
	case "":
		s.home(w, r)
	case "detail":
		head, r.URL.Path = engine.ShiftPath(r.URL.Path)
		id, err := strconv.Atoi(head)
		if err != nil {
			engine.Respond(w, r, http.StatusBadRequest, err)
			return
		}
		log.Println(id)
		//parser functions to be implemented
	default:
		engine.Respond(w, r, http.StatusNotFound, "scoutfish path not found")
	}
}
func (u *Scoutfish) home(w http.ResponseWriter, r *http.Request) {
	jsonString := `{"sub-fen": "r1bqkb1r/pppp1ppp/2n2n2/4p1N1/2B1P3/8/PPPP1PPP/RNBQK2R b KQkq - 0 1"}`

	result, err := helper.NewRunner("../bin/scoutfish/src/scoutfish").Run("scout ../pgn/LumbrasGigaBase-1899.scout", jsonString)
	if err != nil {
		log.Printf("Command execution failed: %v\nOutput: %s", err, result)
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}
	engine.Respond(w, r, http.StatusOK, result)
}
