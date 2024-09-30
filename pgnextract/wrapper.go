package pgnextract

import (
	"engine/helper"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
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

func NewPgnextract(pgn string) *Pgnextract {
	return &Pgnextract{Runner: helper.NewRunner("/app/pgn-extract"), Pgn: fmt.Sprintf("/app/pgn/%s.pgn", pgn), CmdFile: "/app/pgn/pgn-extract-cmd"}
}

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
	command := fmt.Sprintf("%s \"%s\"", cmd, value)
	content := unsafe.Slice(unsafe.StringData(command), len(command))
	err := os.WriteFile(p.CmdFile, content, 0644)
	if err != nil {
		panic(err)
	}
	return true
}

func (p *Pgnextract) PlayerCommand(name string) error {
	if !p.WriteCommand("Player", name) {
		return fmt.Errorf("error generating player command for pgn-extract command file")
	}
	return nil
}

func (p *Pgnextract) WriteMultipleCommands(cmd []string, value []string) bool {

	var fileInput strings.Builder
	for i := 0; i < len(cmd); i++ {
		fileInput.WriteString(fmt.Sprintf("%s \"%s\"\n", cmd[i], value[i]))
	}

	content := unsafe.Slice(unsafe.StringData(fileInput.String()), len(fileInput.String()))

	err := os.WriteFile(p.CmdFile, content, 0644)
	if err != nil {
		panic(err)
	}
	return true
}

func (p *Pgnextract) QueryPlayer(input []byte) ([]helper.PGN, error) {
	var jsonInput PgnInput
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(input, &jsonInput)
	if err != nil {
		log.Println("Error parsing input:", err)
		return nil, err
	}

	playerErr := p.PlayerCommand(jsonInput.Player)
	if playerErr != nil {
		return nil, playerErr
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
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
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
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(input, &jsonInput)
	if err != nil {
		log.Println("Error parsing input:", err)
		return nil, err
	}

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
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(input, &jsonInput)
	if err != nil {
		log.Println("Error parsing input:", err)
		return nil, err
	}

	if !p.WriteCommand("FEN", jsonInput.Fen) {
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
