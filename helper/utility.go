package helper

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
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

// Precompile the reg ex to avoid recompilation
var metadataPattern = regexp.MustCompile(`\[([^\s]+) "([^\"]+)"\]`)

// ParseMetadata parses metadata from a PGN block
func ParseMetadata(game string) map[string]string {
	matches := metadataPattern.FindAllStringSubmatch(game, -1)
	gameMap := make(map[string]string, len(matches)) // Initialize with an expected capacity
	for _, match := range matches {
		gameMap[match[1]] = match[2]
	}
	return gameMap
}

// ParseBoard parses board from a PGN block
func ParseBoard(game string) string {
	lines := strings.Split(game, "\n")
	boardStarted := false
	var board strings.Builder

	for _, line := range lines {
		if !boardStarted && strings.HasPrefix(line, "1.") {
			boardStarted = true
		}
		if boardStarted {
			trimmedLine := strings.TrimSpace(line)
			if trimmedLine == "" || strings.HasPrefix(trimmedLine, "[") {
				continue
			}
			board.WriteString(trimmedLine)
			board.WriteString(" ")
		}
	}
	return strings.TrimSpace(board.String())
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
	var err error

	// Split games based on the occurrence of "[Event "
	games := strings.Split(inputData, "\n\n[Event ")
	if len(games) == 0 {
		err = fmt.Errorf("PGN Parsing failed, no games available")
	}
	if len(games[0]) == 0 {
		games = games[1:]
	}

	var wg sync.WaitGroup
	pgnList := make([]PGN, len(games))
	errors := make(chan error, len(games))
	results := make(chan struct {
		index int
		block PGN
	}, len(games))

	for i, game := range games {
		wg.Add(1)
		go func(i int, game string) {
			defer wg.Done()

			// Ensure the game string starts with [Event
			if !strings.HasPrefix(game, "[Event ") {
				game = "[Event " + game
			}

			// Channels for results of parsing
			metadataChan := make(chan map[string]string)
			boardChan := make(chan string)

			// Both tasks that are executed concurrently starting with go keyword
			// Parse Metadata
			go func() {
				metadataChan <- ParseMetadata(game)
			}()

			// Parse Board
			go func() {
				boardChan <- ParseBoard(game)
			}()

			// Wait for both parsing operations to complete
			gameMap := <-metadataChan
			board := <-boardChan

			// Create PGN struct
			block := CreatePGNStruct(gameMap, board)
			results <- struct {
				index int
				block PGN
			}{index: i, block: block}
		}(i, game)
	}

	// Here the results
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	for result := range results {
		pgnList[result.index] = result.block
	}

	if len(errors) > 0 {
		err = <-errors
	}

	return pgnList, err
}
