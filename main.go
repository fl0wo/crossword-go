package main

import (
	"crossword-go/utils"
	"fmt"
)

func main() {
	// Read words from your JSON file
	data := utils.ReadWords()

	// Extract words for the crossword
	var words []string
	for _, item := range data {
		words = append(words, item.Nome)
		// words = append(words, item.Desc...)
	}

	// shuffle words
	for i := range words {
		j := utils.RandInt(i, len(words))
		words[i], words[j] = words[j], words[i]
	}

	// Create a new crossword puzzle (adjust dimensions as needed)
	puzzle := utils.NewCrossword(15, 15)

	// Generate the puzzle
	success := puzzle.GeneratePuzzle(words)

	if success {
		printCrossWordTerminal(puzzle)

		// Render to PNG
		config := utils.DefaultConfig()
		err := utils.RenderPuzzleToPNG(puzzle, "crossword.png", config)
		if err != nil {
			fmt.Printf("Error rendering puzzle: %v\n", err)
			return
		}

		fmt.Println("Puzzle has been saved to crossword.png")

	} else {
		fmt.Println("Could not generate a valid puzzle with the given words")
	}

}

func printCrossWordTerminal(puzzle *utils.Crossword) {
	// Print the puzzle
	board := puzzle.GetBoard()
	for i := range board {
		for j := range board[i] {
			//fmt.Printf("%c ", board[i][j])
			// always uppercase
			fmt.Printf(string(board[i][j]) + " ")
		}
		fmt.Println()
	}

	// Print word placements
	fmt.Println("\nWord Placements:")
	for _, placement := range puzzle.GetPlacements() {
		direction := "Across"
		if placement.Dir == utils.Vertical {
			direction = "Down"
		}
		fmt.Printf("%s: (%d,%d) %s\n", placement.Word, placement.X, placement.Y, direction)
	}
}
