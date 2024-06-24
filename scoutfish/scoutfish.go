package scoutfish

import (
	"encoding/json"
	"engine"
	"engine/helper"
	"fmt"
	"io/ioutil"
	"net/http"
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
	case "test":
		s.test(w, r)
	default:
		engine.Respond(w, r, http.StatusNotFound, "scoutfish path not found")
	}
}

type Match struct {
	Ofs int   `json:"ofs"`
	Ply []int `json:"ply"`
}

// RunnerOutput represents the JSON structure of the command output
type RunnerOutput struct {
	Moves          int     `json:"moves"`
	MatchCount     int     `json:"match count"`
	MovesPerSecond int     `json:"moves/second"`
	ProcessingTime int     `json:"processing time (ms)"`
	Matches        []Match `json:"matches"`
}

type ScoutfishInput struct {
	SubFen string `json:"sub-fen"`
}

func (u *Scoutfish) test(w http.ResponseWriter, r *http.Request) {
	result, err := helper.NewRunner("./scoutfish").Run("isready")
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	engine.Respond(w, r, http.StatusOK, result)
}

func (u *Scoutfish) home(w http.ResponseWriter, r *http.Request) {
	// jsonString := `{"sub-fen": "r1bqkb1r/pppp1ppp/2n2n2/4p1N1/2B1P3/8/PPPP1PPP/RNBQK2R b KQkq - 0 1"}`

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
		return
	}

	result, err := helper.NewRunner("./scoutfish").Run("scout ../pgn/LumbrasGigaBase-1899.scout", string(body))
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
		return
	}

	var output RunnerOutput

	err = json.Unmarshal([]byte(result), &output)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, err)
	}

	engine.Respond(w, r, http.StatusOK, output)
}
