package helper

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
)

type PGN struct {
	Event     string `json:"event"`
	Site      string `json:"site"`
	Date      string `json:"date"`
	Round     string `json:"round"`
	White     string `json:"white"`
	Black     string `json:"black"`
	Result    string `json:"result"`
	ECO       string `json:"eco"`
	EventDate string `json:"eventDate"`
	PlyCount  string `json:"plyCount"`
	Source    string `json:"source"`
	EventType string `json:"eventType"`
	Board     string `json:"board"`
}

func parsePGNMetadataFromString(pgnData string) (PGN, error) {
	pgn := PGN{}
	scanner := bufio.NewScanner(strings.NewReader(pgnData))
	boardContent := []string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check if this line is part of the metadata
		if len(line) > 2 && line[0] == '[' && line[len(line)-1] == ']' {
			content := line[1 : len(line)-1]
			parts := strings.SplitN(content, " ", 2)
			if len(parts) < 2 {
				continue
			}

			key := parts[0]
			value := strings.Trim(parts[1], "\"") // Remove quotes around value

			switch key {
			case "Event":
				pgn.Event = value
			case "Site":
				pgn.Site = value
			case "Date":
				pgn.Date = value
			case "Round":
				pgn.Round = value
			case "White":
				pgn.White = value
			case "Black":
				pgn.Black = value
			case "Result":
				pgn.Result = value
			case "ECO":
				pgn.ECO = value
			case "EventDate":
				pgn.EventDate = value
			case "PlyCount":
				pgn.PlyCount = value
			case "Source":
				pgn.Source = value
			case "EventType":
				pgn.EventType = value
			}
		} else if line != "" {
			boardContent = append(boardContent, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return PGN{}, err
	}

	pgn.Board = strings.Join(boardContent, " ")
	return pgn, nil
}

func ProcessPGNGames(pgnGamesString string) ([]PGN, error) {
	pgnGames := strings.Split(string(pgnGamesString), "\n\n")

	var wg sync.WaitGroup
	resultChannel := make(chan PGN, len(pgnGames))
	wg.Add(len(pgnGames))

	for _, gameData := range pgnGames {
		go func(data string) {
			defer wg.Done()
			pgn, err := parsePGNMetadataFromString(data)
			if err == nil {
				resultChannel <- pgn
			} else {
				fmt.Println("Error:", err)
				return
			}
		}(gameData)
	}

	wg.Wait()
	close(resultChannel)

	var results []PGN
	for pgn := range resultChannel {
		results = append(results, pgn)
	}
	return results, nil
}
