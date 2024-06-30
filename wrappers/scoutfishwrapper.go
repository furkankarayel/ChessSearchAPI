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
	Ofs int    `json:"ofs"`
	Ply []int  `json:"ply"`
	Pgn string `json:"pgn"`
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

type Scoutfish struct {
	Runner *helper.Runner
	Db     string
	Pgn    string
}

// Our default wrapper initialization that is being used during the development process
func NewDefault() *Scoutfish {
	return &Scoutfish{Runner: helper.NewRunner("/app/scoutfish"), Db: "/app/pgn/LumbrasGigaBase-1899.scout", Pgn: "/app/pgn/LumbrasGigaBase-1899.pgn"}
}

// Wrapper initialization that allows you to choose custom pgn file
func NewWrapper(db string, pgn string) *Scoutfish {
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

func (s *Scoutfish) GetGames(matches []Match) ([]Match, error) {
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

	for i, match := range matches {
		// Seek to the specified offset
		_, err := file.Seek(int64(match.Ofs), 0)
		if err != nil {
			log.Fatal("failed to seek to offset ")
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		var game string
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) > 0 && line[0] == '[' && line[1:7] == "Event " {
				if game != "" {
					break // Second occurrence, start of next game
				}
				game = line // First occurrence
			} else if game != "" {
				game += "\n" + line
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal("error reading PGN file")
			return nil, err
		}
		matches[i].Pgn = game
	}

	return matches, nil

}
