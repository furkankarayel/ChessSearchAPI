package pgnextract

import (
	"strings"
	"testing"
)

var p *Pgnextract

// Initializing the necessary specifications for all tests
func TestMain(m *testing.M) {
	// Initialize the global variable
	p = TestPgnextract("../bin/pgn-extract/src/pgn-extract", "../pgn/LumbrasGigaBase-2020.pgn", "../pgn/pgn-extract-cmd")

	// Run the tests
	m.Run()
}

func TestPgnextractPlayerFirstNameSearch(t *testing.T) {
	playerInput := `{"player": "Alexander"}`

	// Using the player name input as JSON to get the games
	output, err := p.QueryPlayer([]byte(playerInput))
	if err != nil {
		t.Fatalf("Unexpected error when querying player by first name: %v", err)
	}

	// Ensure there are PGN entries returned
	if len(output) == 0 {
		t.Error("Expected to find PGN entries for player with first name 'Alexander', but got none.")
	}

	// Check if the player's first name is in the first PGN entry (could be either White or Black)
	if len(output) > 0 {
		pgn := output[0] // checking the first PGN entry

		if !(strings.Contains(pgn.White, "Alexander") || strings.Contains(pgn.Black, "Alexander")) {
			t.Errorf("Expected player with first name 'Alexander' to be part of the game, but got '%s' vs '%s'.",
				pgn.White, pgn.Black)
		}
	}
}

func TestPgnextractPlayerFullNameSearch(t *testing.T) {
	playerInput := `{"player": "Petkov, Momchil"}`

	// Using the player name input as JSON to get the games
	output, err := p.QueryPlayer([]byte(playerInput))
	if err != nil {
		t.Fatalf("Unexpected error when querying player: %v", err)
	}

	// Ensure there are PGN entries returned
	if len(output) == 0 {
		t.Error("Expected to find PGN entries for player 'Petkov, Momchil', but got none.")
	}

	// Check if the player's name is in the first PGN entry (could be either White or Black)
	if len(output) > 0 {
		pgn := output[0] // checking the first PGN entry

		if !(pgn.White == "Petkov, Momchil" || pgn.Black == "Petkov, Momchil") {
			t.Errorf("Expected player 'Petkov, Momchil' to be part of the game, but got '%s' vs '%s'.",
				pgn.White, pgn.Black)
		}
	}
}

func TestPgnextractTwoSpecificPlayerSearch(t *testing.T) {
	playerInput := `{"players": ["Petkov, Momchil", "Eswaran, Aksithi"]}`

	// Using the player name input as JSON to get the games
	output, err := p.QueryTwoPlayers([]byte(playerInput))
	if err != nil {
		t.Fatalf("Unexpected error when querying players: %v", err)
	}

	// Ensure there are PGN entries returned
	if len(output) == 0 {
		t.Error("Expected to find PGN entries for players 'Petkov, Momchil' and 'Eswaran, Aksithi', but got none.")
	}

	// Check if the player's names are in the first PGN entry (either as White or Black)
	if len(output) > 0 {
		pgn := output[0] // checking the first PGN entry

		if !(pgn.White == "Petkov, Momchil" || pgn.Black == "Petkov, Momchil") ||
			!(pgn.White == "Eswaran, Aksithi" || pgn.Black == "Eswaran, Aksithi") {
			t.Errorf("Expected players 'Petkov, Momchil' and 'Eswaran, Aksithi' to be part of the game, but got '%s' vs '%s'.",
				pgn.White, pgn.Black)
		}
	}
}

// Searching for a player by year of the game
func TestPgnextractPlayerBasedOnYearSearch(t *testing.T) {

	playerInput := `{"player": "Artemiev, Vladislav","year": "2020"}`

	// Using the player name input as json to get the games
	output, err := p.QueryPlayerByYear([]byte(playerInput))
	if err != nil {
		panic(err)
	}

	// Ensure there are PGN entries returned
	if len(output) == 0 {
		t.Error("Expected to find PGN entries for player 'Artemiev, Vladislav' in 2020, but got none.")
	}

	// Check if the player's name is in the first PGN entry (could be either White or Black)
	if len(output) > 0 {
		pgn := output[0] // checking the first PGN entry

		if pgn.White != "Artemiev, Vladislav" && pgn.Black != "Artemiev, Vladislav" {
			t.Errorf("Expected 'Artemiev, Vladislav' to be one of the players, but got '%s' vs '%s'.",
				pgn.White, pgn.Black)
		}
	}

}
