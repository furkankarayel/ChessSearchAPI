package scoutfish

import (
	"engine"
	"engine/wrappers"
	"io/ioutil"
	"net/http"
)

func (u *ScoutfishHandler) test(w http.ResponseWriter, r *http.Request) {
	s := wrappers.NewDefault()

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
	s := wrappers.NewDefault()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	// Using the input json to get the offsets of the games
	scoutraw, err := s.ScoutRaw(body)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	// Using the offsets to get games from the pgn file
	output, err := s.GetGames(scoutraw.Matches)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	engine.Respond(w, r, http.StatusOK, output)
}
