package scoutfish

import (
	"engine"

	"io/ioutil"
	"net/http"
)

func (u *ScoutfishHandler) test(w http.ResponseWriter, r *http.Request) {
	s := DefaultScoutfish()

	result, err := s.IsReady()
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	engine.Respond(w, r, http.StatusOK, string(result))
}

func (u *ScoutfishHandler) home(w http.ResponseWriter, r *http.Request) {
	// jsonString := `{"sub-fen": "r1bqkb1r/pppp1ppp/2n2n2/4p1N1/2B1P3/8/PPPP1PPP/RNBQK2R b KQkq - 0 1"}`

	// Initialize Scoutfish with default settings
	s := DefaultScoutfish()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	// Using the input json to get the offsets of the games
	// Using the input json to query the games
	output, err := s.QueryFen(body)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	engine.Respond(w, r, http.StatusOK, output)
}

// TODO adding new wrapper functions to handler
