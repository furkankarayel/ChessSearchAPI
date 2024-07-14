package wrappers

import (
	"engine/helper"
	"fmt"
	"log"
)

type Pgnextract struct {
	Runner *helper.Runner
	Db     string
	Pgn    string
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

func (p *Pgnextract) QueryFen(input []byte) (*RunnerOutput, error) {
	result, err := p.Runner.Run(fmt.Sprintf("-fen %s %s", string(input), p.Db))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println(result)

	// err = json.Unmarshal([]byte(result), &output)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }
	return nil, nil
}
