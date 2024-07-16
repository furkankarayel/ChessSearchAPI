package wrappers

import (
	"fmt"
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
	Moves     string `json:"Moves"`
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

// ParseMoves parses moves from a PGN block
func ParseMoves(game string) string {
	lines := strings.Split(game, "\n")
	movesStarted := false
	var moves []string

	for _, line := range lines {
		// Identify the start of the moves by checking for move number
		if !movesStarted && strings.HasPrefix(line, "1.") {
			movesStarted = true
		}
		if movesStarted {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "" {
				continue
			}
			if strings.HasPrefix(trimmedLine, "[") {
				break
			}
			moves = append(moves, trimmedLine)
		}
	}
	return strings.Join(moves, " ")
}

// CreatePGNStruct creates a PGN struct from parsed data
func CreatePGNStruct(gameMap map[string]string, moves string) PGN {
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
		Moves:     moves,
	}
}

// ParsePGN parses the entire PGN text and returns a list of PGN structs
func ParsePGN(inputData string) ([]PGN, error) {
	// Split games based on the occurrence of "[Event "
	games := strings.Split(inputData, "\n\n[Event ")

	var pgnList []PGN

	for i, game := range games {

		fmt.Println(fmt.Sprintf("############# %d %s", i, game))
		if i > 0 {
			game = "[Event " + game // Add the [Event tag back to the beginning of each game except the first one
		}
		gameMap := ParseMetadata(game)
		moves := ParseMoves(game)
		pgn := CreatePGNStruct(gameMap, moves)
		pgnList = append(pgnList, pgn)
	}

	return pgnList, nil
}
