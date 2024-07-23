package helper

import (
	"regexp"
	"strings"
)

// PGN represents the structure of a PGN block
type PGN struct {
	Event     string `json:"Event"`
	Site      string `json:"Site"`
	Date      string `json:"Date"`
	Round     string `json:"Round"`
	White     string `json:"White"`
	Black     string `json:"Black"`
	Result    string `json:"Result"`
	ECO       string `json:"ECO"`
	EventDate string `json:"EventDate"`
	PlyCount  string `json:"PlyCount"`
	Source    string `json:"Source"`
	EventType string `json:"EventType"`
	Board     string `json:"Board"`
}

// ParseMetadata parses metadata from a PGN block
func ParseMetadata(game string) map[string]string {
	metadataPattern := regexp.MustCompile(`\[([^\s]+) "([^\"]+)"\]`)
	matches := metadataPattern.FindAllStringSubmatch(game, -1)
	gameMap := make(map[string]string)
	for _, match := range matches {
		gameMap[match[1]] = match[2]
	}
	return gameMap
}

// ParseBoard parses board from a PGN block
func ParseBoard(game string) string {
	lines := strings.Split(game, "\n")
	boardStarted := false
	var board []string

	for _, line := range lines {
		// Identify the start of the board by checking for move number
		if !boardStarted && strings.HasPrefix(line, "1.") {
			boardStarted = true
		}
		if boardStarted {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "" {
				continue
			}
			if strings.HasPrefix(trimmedLine, "[") {
				break
			}
			board = append(board, trimmedLine)
		}
	}
	return strings.Join(board, " ")
}

// CreatePGNStruct creates a PGN struct from parsed data
func CreatePGNStruct(gameMap map[string]string, board string) PGN {
	return PGN{
		Event:     gameMap["Event"],
		Site:      gameMap["Site"],
		Date:      gameMap["Date"],
		Round:     gameMap["Round"],
		White:     gameMap["White"],
		Black:     gameMap["Black"],
		Result:    gameMap["Result"],
		ECO:       gameMap["ECO"],
		EventDate: gameMap["EventDate"],
		PlyCount:  gameMap["PlyCount"],
		Source:    gameMap["Source"],
		EventType: gameMap["EventType"],
		Board:     board,
	}
}

// ParsePGN parses the entire PGN text and returns a list of PGN structs
func ParsePGN(inputData string) ([]PGN, error) {
	// Split games based on the occurrence of "[Event "
	games := strings.Split(inputData, "\n\n[Event ")

	var pgnList []PGN

	for i, game := range games {
		if i > 0 {
			game = "[Event " + game // Add the [Event tag back to the beginning of each game except the first one
		}
		gameMap := ParseMetadata(game)
		board := ParseBoard(game)
		pgn := CreatePGNStruct(gameMap, board)
		pgnList = append(pgnList, pgn)
	}

	return pgnList, nil
}
