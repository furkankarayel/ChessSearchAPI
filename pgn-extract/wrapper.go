package pgnextract

import (
	"encoding/json"
	"engine/helper"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type Pgnextract struct {
	Runner  *helper.Runner
	Pgn     string
	CmdFile string
}

type PgnInput struct {
	Fen     string
	Player  string
	Players []string
	Year    string
}

// Our default wrapper initialization that is being used during the development process
func DefaultPgnextract() *Pgnextract {
	return &Pgnextract{Runner: helper.NewRunner("/app/pgn-extract"), Pgn: "/app/pgn/LumbrasGigaBase-1899.pgn"}
}

// Wrapper initialization that allows you to choose custom pgn file
func NewPgnextract(db string, pgn string) *Pgnextract {
	return &Pgnextract{Runner: helper.NewRunner("/app/pgn-extract"), Pgn: fmt.Sprintf("/app/pgn/%s.pgn", pgn)}
}

// Wrapper initialization for tests that allows you to pass binary and pgn file
func TestPgnextract(binary string, pgn string, cmdFile string) *Pgnextract {
	return &Pgnextract{Runner: helper.NewRunner(binary), Pgn: pgn, CmdFile: cmdFile}
}

func (p *Pgnextract) IsReady() ([]byte, error) {
	result, err := p.Runner.Run("/app/pgn/test.pgn")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return result, nil
}

func (p *Pgnextract) WriteCommand(cmd string, value string) bool {

	content := []byte(fmt.Sprintf("%s \"%s\"", cmd, value))
	err := os.WriteFile(p.CmdFile, content, 0644)
	if err != nil {
		panic(err)
	}
	return true
}

func (p *Pgnextract) WriteMultipleCommands(cmd []string, value []string) bool {

	fileInput := ""
	for i := 0; i < len(cmd); i++ {
		fileInput += fmt.Sprintf("%s \"%s\"\n", cmd[i], value[i])
	}

	content := []byte(fileInput)

	err := os.WriteFile(p.CmdFile, content, 0644)
	if err != nil {
		panic(err)
	}
	return true
}

// Querying Player name Input
func (p *Pgnextract) QueryPlayer(input []byte) ([]helper.PGN, error) {
	var jsonInput PgnInput
	err := json.Unmarshal(input, &jsonInput)
	if err != nil {
		log.Println("Error parsing input:", err)
		return nil, err
	}

	if !p.WriteCommand("Player", jsonInput.Player) {
		log.Println("Error generating pgn-extract command file:", err)
		return nil, err
	}

	// Generate a UUID for the temporary output file name
	tempOutputFilename := uuid.New().String()

	// Ensure the temporary file is removed after usage
	defer func() {
		if err := os.Remove(tempOutputFilename); err != nil {
			log.Printf("Failed to remove temporary file: %v", err)
		}
	}()

	// Pgn-extract is processing commands out of a file and saves the results as file lastOutput
	_, err = p.Runner.Run(fmt.Sprintf("-t %s", p.CmdFile), p.Pgn, "--output", tempOutputFilename)
	if err != nil {
		log.Println("Error finding player:", err)
		return nil, err
	}

	// Clean approach to process the output of the command by an output file instead of terminal output
	// as it adds an extra line of match players summary in terminal
	catRunner := helper.NewRunner("cat")
	readOutput, err := catRunner.Run(tempOutputFilename)
	if err != nil {
		log.Println("Error finding player:", err)
		return nil, err
	}

	pgnList, err := helper.ParsePGN(string(readOutput))
	if err != nil {
		fmt.Println("Error finding player:", err)
		return nil, err
	}

	return pgnList, nil
}

// Querying Player name Input
func (p *Pgnextract) QueryTwoPlayers(input []byte) ([]helper.PGN, error) {
	var jsonInput PgnInput
	err := json.Unmarshal(input, &jsonInput)
	if err != nil {
		log.Println("Error parsing input:", err)
		return nil, err
	}

	// Writing commands to cover both color combinations for each player given
	// Results in all games both have played against each other having different colors
	if !p.WriteMultipleCommands([]string{"Black", "White", "White", "Black"}, []string{string(jsonInput.Players[0]), string(jsonInput.Players[1]), string(jsonInput.Players[0]), string(jsonInput.Players[1])}) {
		log.Println("Error generating pgn-extract command file:", err)
		return nil, err
	}

	// Generate a UUID for the temporary output file name
	tempOutputFilename := uuid.New().String()

	// Ensure the temporary file is removed after usage
	defer func() {
		if err := os.Remove(tempOutputFilename); err != nil {
			log.Printf("Failed to remove temporary file: %v", err)
		}
	}()

	// Pgn-extract is processing commands out of a file and saves the results as file lastOutput
	_, err = p.Runner.Run(fmt.Sprintf("-t %s", p.CmdFile), p.Pgn, "--output", tempOutputFilename)
	if err != nil {
		log.Println("Error finding player:", err)
		return nil, err
	}

	// Clean approach to process the output of the command by an output file instead of terminal output
	// as it adds an extra line of match players summary in terminal
	catRunner := helper.NewRunner("cat")
	readOutput, err := catRunner.Run(tempOutputFilename)
	if err != nil {
		log.Println("Error finding player:", err)
		return nil, err
	}

	pgnList, err := helper.ParsePGN(string(readOutput))
	if err != nil {
		fmt.Println("Error finding player:", err)
		return nil, err
	}

	return pgnList, nil
}

// Querying Player name Input
func (p *Pgnextract) QueryPlayerByYear(input []byte) ([]helper.PGN, error) {
	var jsonInput PgnInput
	err := json.Unmarshal(input, &jsonInput)
	if err != nil {
		log.Println("Error parsing input:", err)
		return nil, err
	}

	fmt.Print(jsonInput.Player)
	fmt.Print(jsonInput.Year)
	// Writing commands to query player and date
	if !p.WriteMultipleCommands([]string{"Player", "Date"}, []string{string(jsonInput.Player), string(jsonInput.Year)}) {
		log.Println("Error generating pgn-extract command file:", err)
		return nil, err
	}

	// Generate a UUID for the temporary output file name
	tempOutputFilename := uuid.New().String()

	// Ensure the temporary file is removed after usage
	defer func() {
		if err := os.Remove(tempOutputFilename); err != nil {
			log.Printf("Failed to remove temporary file: %v", err)
		}
	}()

	// Pgn-extract is processing commands out of a file and saves the results as file lastOutput
	_, err = p.Runner.Run(fmt.Sprintf("-t %s", p.CmdFile), p.Pgn, "--output", tempOutputFilename)
	if err != nil {
		log.Println("Error finding player:", err)
		return nil, err
	}

	// Clean approach to process the output of the command by an output file instead of terminal output
	// as it adds an extra line of match players summary in terminal
	catRunner := helper.NewRunner("cat")
	readOutput, err := catRunner.Run(tempOutputFilename)
	if err != nil {
		log.Println("Error finding player:", err)
		return nil, err
	}

	pgnList, err := helper.ParsePGN(string(readOutput))
	if err != nil {
		fmt.Println("Error finding player:", err)
		return nil, err
	}

	return pgnList, nil
}

// Querying FEN Input
func (p *Pgnextract) QueryFen(input []byte) ([]helper.PGN, error) {
	var jsonInput PgnInput
	err := json.Unmarshal(input, &jsonInput)
	if err != nil {
		log.Println("Error parsing input:", err)
		return nil, err
	}

	if !p.WriteCommand("Player", jsonInput.Fen) {
		log.Println("Error generating pgn-extract command file:", err)
		return nil, err
	}

	// Generate a UUID for the temporary output file name
	tempOutputFilename := uuid.New().String()

	// Ensure the temporary file is removed after usage
	defer func() {
		if err := os.Remove(tempOutputFilename); err != nil {
			log.Printf("Failed to remove temporary file: %v", err)
		}
	}()

	// Pgn-extract is processing commands out of a file and saves the results as file lastOutput
	_, err = p.Runner.Run(fmt.Sprintf("-t %s ", p.CmdFile), p.Pgn, "--output", tempOutputFilename)
	if err != nil {
		log.Println("Error finding FEN:", err)
		return nil, err
	}

	// Clean approach to process the output of the command by an output file instead of terminal output
	// as it adds an extra line of match players summary in terminal
	catRunner := helper.NewRunner("cat")
	readOutput, err := catRunner.Run("lastOutput")
	if err != nil {
		log.Println("Error finding FEN:", err)
		return nil, err
	}

	pgnList, err := helper.ParsePGN(string(readOutput))
	if err != nil {
		fmt.Println("Error parsing PGN:", err)
		return nil, err
	}

	return pgnList, nil
}

//TODO Search for two specific players

//TODO Search for two specific players where only 1 player won

// TODO criteria builder which is writing the commands into a file, that is going to be used for executing queries
