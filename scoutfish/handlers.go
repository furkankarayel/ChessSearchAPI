package scoutfish

import (
	"encoding/json"
	"engine"
	"fmt"
	"unsafe"

	"net/http"
)

type FenRequest struct {
	SubFen string `json:"sub-fen"`
}

type ArbitraryPositionRequest struct {
	WhiteMove string `json:"white-move"`
}

func (u *ScoutfishHandler) searchByFen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		engine.Respond(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var fenReq FenRequest
	err := json.NewDecoder(r.Body).Decode(&fenReq)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, fmt.Sprintf("Failed to decode FEN request: %s", err))
		return
	}

	byteInput := unsafe.Slice(unsafe.StringData(fenReq.SubFen), len(fenReq.SubFen))

	output, err := u.scoutfish.QueryFen(byteInput)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, fmt.Sprintf("Error querying FEN: %s", err))
		return
	}

	if len(output) == 0 {
		engine.Respond(w, r, http.StatusNotFound, "No results found for the given FEN position")
		return
	}

	engine.Respond(w, r, http.StatusOK, output)
}

func (u *ScoutfishHandler) searchArbitraryPosition(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		engine.Respond(w, r, http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var posReq ArbitraryPositionRequest
	err := json.NewDecoder(r.Body).Decode(&posReq)
	if err != nil {
		engine.Respond(w, r, http.StatusBadRequest, fmt.Sprintf("Failed to decode position request: %s", err))
		return
	}

	byteInput := unsafe.Slice(unsafe.StringData(posReq.WhiteMove), len(posReq.WhiteMove))

	output, err := u.scoutfish.QueryFen(byteInput)
	if err != nil {
		engine.Respond(w, r, http.StatusInternalServerError, fmt.Sprintf("Error querying position: %s", err))
		return
	}

	if len(output) == 1 {
		engine.Respond(w, r, http.StatusNotFound, "No results found for the given arbitrary position")
		return
	}

	engine.Respond(w, r, http.StatusOK, output)
}
