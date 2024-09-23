package scoutfish

import (
	"engine/helper"
	"testing"
	"unsafe"
)

var s *Scoutfish

// Initializing the necessary specifications for all tests
func TestMain(m *testing.M) {
	// Initialize the global variable
	s = TestScoutfish("../bin/scoutfish/src/scoutfish", "../pgn/LumbrasGigaBase-2020.scout", "../pgn/LumbrasGigaBase-2020.pgn")

	// Run the tests
	m.Run()
}

func ProcessQuery(input string) []helper.PGN {

	byteInput := unsafe.Slice(unsafe.StringData(input), len(input))

	// Using the input json to query the games
	output, err := s.QueryFen(byteInput)
	if err != nil {
		panic(err)
	}
	return output
}

func TestScoutfishSearchFENPosition(t *testing.T) {
	// A specific FEN input to search for in the database
	input := `{"sub-fen": "r1bqkb1r/pppp1ppp/2n2n2/4p1N1/2B1P3/8/PPPP1PPP/RNBQK2R b KQkq - 0 1"}`

	byteInput := unsafe.Slice(unsafe.StringData(input), len(input))
	// Call the function to process the query
	output, err := s.QueryFen(byteInput)
	if err != nil {
		panic(err)
	}

	// Ensure there is output; length of 1 usually indicates no result
	if len(output) == 0 {
		t.Error("Expected results for the FEN position search, but got none.")
	}
}

// Search for an arbitrary chess position where the pieces can be anywhere on the board
func Test_Scoutfish_Search_Arbitrary_Position(t *testing.T) {

	// An arbitrary chess position of the white player that is being searched for in the database
	input := `{"white-move": "e5"}`

	byteInput := unsafe.Slice(unsafe.StringData(input), len(input))

	// Call the function to process the query
	output, err := s.QueryFen(byteInput)
	if err != nil {
		panic(err)
	}

	// No result always gives 1 Byte back
	if len(output) == 1 {
		t.Error("Scoutfish Search Arbitrary Position: not found")
	}

}
