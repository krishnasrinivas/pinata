/*
Copyright © 2020 Anand Babu Periasamy https://twitter.com/abperiasamy

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/abperiasamy/chess"
)

// Initalize a new game.
func newGame() *chess.Game {
	// Use human friendly short Algebraic notation (like e4, e5)
	return chess.NewGame(chess.UseNotation(chess.AlgebraicNotation{}))
}

// Start a game from a PGN file
func loadPGN(filename string) *chess.Game {
	pgnDat, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Unable to open", gConsole.Bold(gConsole.Yellow(filename)))
		return nil
	}

	pgn, err := chess.PGN(strings.NewReader(string(pgnDat)))
	if err != nil {
		fmt.Println(gConsole.Bold(gConsole.Yellow(filename)), "is not a valid PGN file.")
		return nil
	}
	return chess.NewGame(pgn)
}

// Save the game to a PGN file
func savePGN(game *chess.Game, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Unable to create", gConsole.Bold(gConsole.Yellow(filename)))
		return err
	}
	defer file.Close()

	// Generate PGN content.
	curTime := time.Now()
	curDate := fmt.Sprintf("%d-%02d-%02d", curTime.Year(), curTime.Month(), curTime.Day())
	game.AddTagPair("Date", curDate)
	game.AddTagPair("Result", game.Outcome().String())
	if humanColor() == chess.White {
		game.AddTagPair("White", "Human")
		game.AddTagPair("Black", gEngineBinary)
	} else {
		game.AddTagPair("White", gEngineBinary)
		game.AddTagPair("Black", "Human")
	}

	// Save the engine name.
	_, err = file.WriteString(game.String() + "\n")
	if err != nil {
		fmt.Println("Unable to save the game to", gConsole.Bold(gConsole.Yellow(filename)))
		return err
	}

	return nil // Success
}

func drawBoard(game *chess.Game) {
	if !gVisual {
		return // playing blind
	}

	if gHumanIsBlack { // Rotate the board, black facing the human.
		fmt.Print(game.Position().Board().DrawForBlack())
	} else {
		fmt.Print(game.Position().Board().Draw())
	}
}

func isGameOver(game *chess.Game) bool {
	switch game.Outcome() {
	case chess.NoOutcome:
		return false
	case chess.Draw:
		fmt.Println(gConsole.Bold(gConsole.Yellow("Game draw!!")))
	case chess.WhiteWon:
		fmt.Println(gConsole.Bold(gConsole.Yellow("White won the game!!")))
	case chess.BlackWon:
		fmt.Println(gConsole.Bold(gConsole.Yellow("Black won the game!!")))
	default:
		panic(game.Outcome()) // should never happen
	}
	return true // The end.
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func humanColor() chess.Color {
	if gHumanIsBlack {
		return chess.Black
	}
	return chess.White
}

// Readline completion of all the valid moves left.
func validMovesConstructor(game *chess.Game) func(string) []string {
	return func(string) (moves []string) {
		for _, move := range game.Position().ValidMoves() {
			moveSAN := chess.Encoder.Encode(chess.AlgebraicNotation{}, game.Position(), move)
			moves = append(moves, moveSAN)
		}
		return moves
	}
}

// Readline completion of all the valid moves left.
func validMoves(game *chess.Game) (moves string) {
	for _, move := range game.Position().ValidMoves() {
		moves += " " + chess.Encoder.Encode(chess.AlgebraicNotation{}, game.Position(), move)
	}
	return moves
}
