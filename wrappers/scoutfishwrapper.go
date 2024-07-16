package wrappers

import (
	"bufio"
	"encoding/json"
	"engine/helper"
	"fmt"
	"log"
	"os"
)

type Match struct {
	Ofs   int    `json:"ofs"`
	Ply   []int  `json:"ply"`
	Board string `json:"board"`
}

// RunnerOutput represents the JSON structure of the command output
type RunnerOutput struct {
	Moves          int     `json:"moves"`
	MatchCount     int     `json:"match count"`
	MovesPerSecond int     `json:"moves/second"`
	ProcessingTime int     `json:"processing time (ms)"`
	Matches        []Match `json:"matches"`
}

type Scoutfish struct {
	Runner *helper.Runner
	Db     string
	Pgn    string
}

// Our default wrapper initialization that is being used during the development process
func DefaultScoutfish() *Scoutfish {
	return &Scoutfish{Runner: helper.NewRunner("/app/scoutfish"), Db: "/app/pgn/LumbrasGigaBase-1899.scout", Pgn: "/app/pgn/LumbrasGigaBase-1899.pgn"}
}

// Wrapper initialization that allows you to choose custom pgn file
func NewScoutfish(db string, pgn string) *Scoutfish {
	return &Scoutfish{Runner: helper.NewRunner("/app/scoutfish"), Db: fmt.Sprintf("/app/pgn/%s.pgn", db), Pgn: fmt.Sprintf("/app/pgn/%s.pgn", pgn)}
}

func (s *Scoutfish) IsReady() ([]byte, error) {
	result, err := s.Runner.Run("isready")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Scoutfish) ScoutRaw(input []byte) (*RunnerOutput, error) {
	result, err := s.Runner.Run(fmt.Sprintf("scout %s %s", s.Db, string(input)))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var output RunnerOutput

	err = json.Unmarshal([]byte(result), &output)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &output, nil

}

func (s *Scoutfish) GetGames(matches []Match) ([]PGN, error) {
	if len(s.Pgn) == 0 {
		log.Fatal("PGN File not found..")
		return nil, nil
	}

	file, err := os.Open(s.Pgn)
	if err != nil {
		log.Fatal("Couldn't open PGN file..")
		return nil, err
	}
	defer file.Close()

	output := ""
	for _, match := range matches {
		// Seek to the specified offset
		_, err := file.Seek(int64(match.Ofs), 0)
		if err != nil {
			log.Fatal("failed to seek to offset: ", err)
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		var game string
		isGameStarted := false
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) > 0 && line[0] == '[' && line[1:7] == "Event " {
				if isGameStarted {
					break // Second occurrence, start of next game
				}
				isGameStarted = true // First occurrence
				game = line
			} else if isGameStarted {
				game += "\n" + line
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal("error reading PGN file: ", err)
			return nil, err
		}

		output += game + "\n\n" // Separate games with a double newline
	}

	pgnList, err := ParsePGN(output)
	if err != nil {
		fmt.Println("Error parsing PGN:", err)
		return nil, err
	}

	return pgnList, nil

}
