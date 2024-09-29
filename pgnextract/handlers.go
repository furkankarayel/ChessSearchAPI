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

	// Using the player name input as JSON to get the games
	output, err := u.pgnextract.QueryPlayer([]byte(playerInput))
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, fmt.Sprintf("Unexpected error when querying player by first name: %s", err))
	}

	// Ensure there are PGN entries returned
	if len(output) == 0 {
		engine.Respond(w, r, http.StatusBadRequest, fmt.Sprintf("Expected to find PGN entries for player with first name '%s', but got none.", playerInput))
	}
	engine.Respond(w, r, http.StatusOK, output)
}
