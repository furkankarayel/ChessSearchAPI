package pgnextract

import (
	"encoding/json"
	"engine"
	"fmt"
	"net/http"
)

type PlayerRequest struct {
	Name string `json:"name"`
}

type TwoPlayersRequest struct {
	Players []string `json:"players"`
}

type PlayerYearRequest struct {
	Player string `json:"player"`
	Year   int    `json:"year"`
}

func (u *PgnextracthHandler) searchPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		engine.Respond(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var playerReq PlayerRequest
	err := json.NewDecoder(r.Body).Decode(&playerReq)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, "Internal server error")
		return
	}

	playerName := playerReq.Name
	playerInput := fmt.Sprintf(`{"player": "%s"}`, playerName)

	output, err := u.pgnextract.QueryPlayer([]byte(playerInput))
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, fmt.Sprintf("Unexpected error when querying player by first name: %s", err))
	}

	if len(output) == 0 {
		engine.Respond(w, r, http.StatusBadRequest, fmt.Sprintf("Expected to find PGN entries for player with first name '%s', but got none.", playerInput))
	}
	engine.Respond(w, r, http.StatusOK, output)
}

func (u *PgnextracthHandler) searchTwoPlayers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		engine.Respond(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var playersReq TwoPlayersRequest
	err := json.NewDecoder(r.Body).Decode(&playersReq)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, "Internal server error")
		return
	}

	if len(playersReq.Players) != 2 {
		engine.Respond(w, r, http.StatusBadRequest, "Exactly two players must be provided")
		return
	}

	playerInput := fmt.Sprintf(`{"players": ["%s", "%s"]}`, playersReq.Players[0], playersReq.Players[1])
	output, err := u.pgnextract.QueryTwoPlayers([]byte(playerInput))
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, fmt.Sprintf("Unexpected error when querying players: %s", err))
		return
	}

	if len(output) == 0 {
		engine.Respond(w, r, http.StatusNotFound, fmt.Sprintf("Expected to find PGN entries for players '%s' and '%s', but got none.", playersReq.Players[0], playersReq.Players[1]))
		return
	}

	pgn := output[0]
	if !(pgn.White == playersReq.Players[0] || pgn.Black == playersReq.Players[0]) ||
		!(pgn.White == playersReq.Players[1] || pgn.Black == playersReq.Players[1]) {
		engine.Respond(w, r, http.StatusBadRequest, fmt.Sprintf("Expected players '%s' and '%s' to be part of the game, but got '%s' vs '%s'.", playersReq.Players[0], playersReq.Players[1], pgn.White, pgn.Black))
		return
	}

	engine.Respond(w, r, http.StatusOK, output)
}

func (u *PgnextracthHandler) searchPlayerByYear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		engine.Respond(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var playerYearReq PlayerYearRequest
	err := json.NewDecoder(r.Body).Decode(&playerYearReq)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, fmt.Sprintf("Failed to decode player year request: %s", err))
		return
	}

	playerInput := fmt.Sprintf(`{"player": "%s", "year": %d}`, playerYearReq.Player, playerYearReq.Year)
	output, err := u.pgnextract.QueryPlayerByYear([]byte(playerInput))
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, fmt.Sprintf("Unexpected error when querying player by year: %s", err))
		return
	}

	if len(output) == 0 {
		engine.Respond(w, r, http.StatusNotFound, fmt.Sprintf("Expected to find PGN entries for player '%s' in year %d, but got none.", playerYearReq.Player, playerYearReq.Year))
		return
	}

	pgn := output[0]
	if pgn.White != playerYearReq.Player && pgn.Black != playerYearReq.Player {
		engine.Respond(w, r, http.StatusBadRequest, fmt.Sprintf("Expected player '%s' to be part of the game in year %d, but got '%s' vs '%s'.", playerYearReq.Player, playerYearReq.Year, pgn.White, pgn.Black))
		return
	}

	engine.Respond(w, r, http.StatusOK, output)
}
