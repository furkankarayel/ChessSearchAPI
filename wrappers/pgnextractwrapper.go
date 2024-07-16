package wrappers

import (
	"encoding/json"
	"engine/helper"
	"fmt"
	"log"
)

type Pgnextract struct {
	Runner *helper.Runner
	Db     string
	Pgn    string
}

type PgnInput struct {
	Fen string
}

// Our default wrapper initialization that is being used during the development process
func DefaultPgnextract() *Pgnextract {
	return &Pgnextract{Runner: helper.NewRunner("/app/pgn-extract"), Db: "/app/pgn/LumbrasGigaBase-1899.scout", Pgn: "/app/pgn/LumbrasGigaBase-1899.pgn"}
}

// Wrapper initialization that allows you to choose custom pgn file
func NewPgnextract(db string, pgn string) *Pgnextract {
	return &Pgnextract{Runner: helper.NewRunner("/app/pgn-extract"), Db: fmt.Sprintf("/app/pgn/%s.pgn", db), Pgn: fmt.Sprintf("/app/pgn/%s.pgn", pgn)}
}

func (p *Pgnextract) IsReady() ([]byte, error) {
	result, err := p.Runner.Run("/app/pgn/test.pgn")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return result, nil
}

func (p *Pgnextract) QueryFen(input []byte) ([]PGN, error) {
	var jsonInput PgnInput
	err := json.Unmarshal(input, &jsonInput)
	if err != nil {
		log.Println("Error parsing input:", err)
		return nil, err
	}

	_, err = p.Runner.Run("-t /app/pgn/pgn-test-command", p.Pgn, "--output", "lastOutput")
	if err != nil {
		log.Println("Error finding FEN:", err)
		return nil, err
	}

	p.Runner.Executable = "cat"
	readOutput, err := p.Runner.Run("lastOutput")
	if err != nil {
		log.Println("Error finding FEN:", err)
		return nil, err
	}

	pgnList, err := ParsePGN(string(readOutput))
	if err != nil {
		fmt.Println("Error parsing PGN:", err)
		return nil, err
	}

	return pgnList, nil
}
