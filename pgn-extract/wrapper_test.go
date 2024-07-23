package pgnextract

import (
	"fmt"
	"testing"
)

var p *Pgnextract

// Initializing the necessary specifications for all tests
func TestMain(m *testing.M) {
	// Initialize the global variable
	p = TestPgnextract("../bin/pgn-extract/src/pgn-extract", "../pgn/LumbrasGigaBase-1899.pgn", "../pgn/pgn-extract-cmd")

	// Run the tests
	m.Run()
}

// Searching for a player by first name
func Test_Pgnextract_Player_FirstName_Search(t *testing.T) {

	playerInput := `{"player": "Alexander"}`

	// Using the player name input as json to get the games
	output, err := p.QueryPlayer([]byte(playerInput))
	if err != nil {
		panic(err)
	}

	// No result always gives 1 Byte back
	if len(output) == 1 {
		t.Error("Player search by first name: not found")
	}
}

// Searching for a player by full name name
func Test_Pgnextract_Player_FullName_Search(t *testing.T) {

	playerInput := `{"player": "Cordel, Oskar"}`

	// Using the player name input as json to get the games
	output, err := p.QueryPlayer([]byte(playerInput))
	if err != nil {
		panic(err)
	}

	// No result always gives 1 Byte back
	if len(output) == 1 {
		t.Error("Player search by full name: not found")
	}

}

// Searching for a player by full name name
func Test_Pgnextract_Two_Specific_Player_Search(t *testing.T) {

	playerInput := `{"players": ["Suhle, Berthold","Cordel, Oskar"]}`

	// Using the player name input as json to get the games
	output, err := p.QueryTwoPlayers([]byte(playerInput))
	if err != nil {
		panic(err)
	}

	// No result always gives 1 Byte back
	if len(output) == 1 {
		t.Error("Specific players search by full name: not found")
	}

}

// Searching for a player by year of the game
func Test_Pgnextract_Player_Based_On_Year_Search(t *testing.T) {

	playerInput := `{"player": "Magnus Carlsen","year": "2020"}`

	// Using the player name input as json to get the games
	output, err := p.QueryPlayerByYear([]byte(playerInput))
	if err != nil {
		panic(err)
	}

	// No result always gives 1 Byte back
	if len(output) == 1 {
		t.Error("Player search by full name and year: not found")
	}

	fmt.Print(output)
}

// TODO tests need to be specified, they're not good enough to check if valid information are given back.
// Currently they only check if any response bigger than 1 Byte is given back. Valid responses were checked manually
