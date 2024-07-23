package scoutfish

import (
	"testing"
)

var s *Scoutfish

// Initializing the necessary specifications for all tests
func TestMain(m *testing.M) {
	// Initialize the global variable
	s = TestScoutfish("../bin/scoutfish/src/scoutfish", "../pgn/LumbrasGigaBase-1899.scout", "../pgn/LumbrasGigaBase-1899.pgn")

	// Run the tests
	m.Run()
}

// Searching for specific fen positions
func Test_Scoutfish_Search_FEN_Position(t *testing.T) {
	// A specific FEN Input that is being searched for in the database
	input := `{"sub-fen": "r1bqkb1r/pppp1ppp/2n2n2/4p1N1/2B1P3/8/PPPP1PPP/RNBQK2R b KQkq - 0 1"}`
	byteInput := []byte(input)

	// Using the input json to query the games
	output, err := s.QueryFen(byteInput)
	if err != nil {
		panic(err)
	}

	// No result always gives 1 Byte back
	if len(output) == 1 {
		t.Error("Scoutfish specific FEN Position: not found")
	}

}

// Search for an arbitrary chess position where the pieces can be anywhere on the board
func Test_Scoutfish_Search_Arbitrary_Position(t *testing.T) {

	// An arbitrary chess position of the white player that is being searched for in the database
	input := `{"white-move": "e5"}`
	byteInput := []byte(input)

	// Using the input json to query the games
	output, err := s.QueryFen(byteInput)
	if err != nil {
		panic(err)
	}

	// No result always gives 1 Byte back
	if len(output) == 1 {
		t.Error("Scoutfish Search Arbitrary Position: not found")
	}

}
