package pgnextract

import (
	"engine"
	"io/ioutil"
	"net/http"
)

// Specific approach to run a test for pgn-extract as it doesn't have any built in check functionality
func (u *PgnextracthHandler) test(w http.ResponseWriter, r *http.Request) {
	p := DefaultPgnextract()

	result, err := p.IsReady()
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}
	// do nothing, solution to declared and not used error
	_ = result

	engine.Respond(w, r, http.StatusOK, "readyok")
}

func (u *PgnextracthHandler) home(w http.ResponseWriter, r *http.Request) {
	// jsonString := `{"sub-fen": "r1bqkb1r/pppp1ppp/2n2n2/4p1N1/2B1P3/8/PPPP1PPP/RNBQK2R b KQkq - 0 1"}`

	// Initialize Scoutfish with default settings
	p := DefaultPgnextract()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	// Using the input json to get the offsets of the games
	output, err := p.QueryFen(body)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	engine.Respond(w, r, http.StatusOK, output)
}

// TODO adding new wrapper functions to handler
