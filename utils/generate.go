package utils

import (
	"math/rand"
	"time"
)

// Direction represents the orientation of a word
type Direction int

const (
	Horizontal Direction = 0
	Vertical   Direction = 1
)

// Position represents a starting position and direction for a word
type Position struct {
	X, Y int
	Dir  Direction
}

// WordPlacement represents a placed word and its metadata
type WordPlacement struct {
	X, Y   int
	Dir    Direction
	Length int
	Word   string
}

// Crossword represents the crossword puzzle
type Crossword struct {
	board      [][]rune
	hWords     [][]int // tracks horizontal words
	vWords     [][]int // tracks vertical words
	width      int
	height     int
	usedWords  map[string]bool
	placements []WordPlacement
	hCount     int
	vCount     int
}

// NewCrossword creates a new crossword puzzle with given dimensions
func NewCrossword(width, height int) *Crossword {
	c := &Crossword{
		width:     width,
		height:    height,
		usedWords: make(map[string]bool),
	}

	// Initialize the board
	c.board = make([][]rune, height)
	c.hWords = make([][]int, height)
	c.vWords = make([][]int, height)

	for i := 0; i < height; i++ {
		c.board[i] = make([]rune, width)
		c.hWords[i] = make([]int, width)
		c.vWords[i] = make([]int, width)
		for j := 0; j < width; j++ {
			c.board[i][j] = ' '
		}
	}

	return c
}

// isValidPosition checks if the given coordinates are within the board
func (c *Crossword) isValidPosition(x, y int) bool {
	return x >= 0 && x < c.height && y >= 0 && y < c.width
}

// canBePlaced checks if a word can be placed at the given position
func (c *Crossword) canBePlaced(word string, x, y int, dir Direction) int {
	intersections := 0

	if dir == Horizontal {
		// Check horizontal placement
		for j := 0; j < len(word); j++ {
			x1, y1 := x, y+j

			if !c.isValidPosition(x1, y1) {
				return -1
			}

			// Check if space is empty or matches letter
			if c.board[x1][y1] != ' ' && c.board[x1][y1] != rune(word[j]) {
				return -1
			}

			// Check adjacent words
			if c.isValidPosition(x1-1, y1) && c.hWords[x1-1][y1] > 0 {
				return -1
			}
			if c.isValidPosition(x1+1, y1) && c.hWords[x1+1][y1] > 0 {
				return -1
			}

			if c.board[x1][y1] == rune(word[j]) {
				intersections++
			}
		}
	} else {
		// Check vertical placement
		for j := 0; j < len(word); j++ {
			x1, y1 := x+j, y

			if !c.isValidPosition(x1, y1) {
				return -1
			}

			if c.board[x1][y1] != ' ' && c.board[x1][y1] != rune(word[j]) {
				return -1
			}

			if c.isValidPosition(x1, y1-1) && c.vWords[x1][y1-1] > 0 {
				return -1
			}
			if c.isValidPosition(x1, y1+1) && c.vWords[x1][y1+1] > 0 {
				return -1
			}

			if c.board[x1][y1] == rune(word[j]) {
				intersections++
			}
		}
	}

	// Check spaces before and after word
	if dir == Horizontal {
		if c.isValidPosition(x, y-1) && c.board[x][y-1] != ' ' && c.board[x][y-1] != '*' {
			return -1
		}
		if c.isValidPosition(x, y+len(word)) && c.board[x][y+len(word)] != ' ' && c.board[x][y+len(word)] != '*' {
			return -1
		}
	} else {
		if c.isValidPosition(x-1, y) && c.board[x-1][y] != ' ' && c.board[x-1][y] != '*' {
			return -1
		}
		if c.isValidPosition(x+len(word), y) && c.board[x+len(word)][y] != ' ' && c.board[x+len(word)][y] != '*' {
			return -1
		}
	}

	return intersections
}

// putWord places a word on the board
func (c *Crossword) putWord(word string, x, y int, dir Direction) {
	if c.usedWords[word] {
		return
	}

	value := 0
	if dir == Horizontal {
		c.hCount++
		value = c.hCount
	} else {
		c.vCount++
		value = c.vCount
	}

	c.usedWords[word] = true
	c.placements = append(c.placements, WordPlacement{
		X:      x,
		Y:      y,
		Dir:    dir,
		Length: len(word),
		Word:   word,
	})

	for i := 0; i < len(word); i++ {
		var x1, y1 int
		if dir == Horizontal {
			x1, y1 = x, y+i
		} else {
			x1, y1 = x+i, y
		}

		c.board[x1][y1] = rune(word[i])
		if dir == Horizontal {
			c.hWords[x1][y1] = value
		} else {
			c.vWords[x1][y1] = value
		}
	}

	// Place blocking characters
	if dir == Horizontal {
		if c.isValidPosition(x, y-1) {
			c.board[x][y-1] = '*'
		}
		if c.isValidPosition(x, y+len(word)) {
			c.board[x][y+len(word)] = '*'
		}
	} else {
		if c.isValidPosition(x-1, y) {
			c.board[x-1][y] = '*'
		}
		if c.isValidPosition(x+len(word), y) {
			c.board[x+len(word)][y] = '*'
		}
	}
}

// findBestPosition finds the best position for a word
func (c *Crossword) findBestPosition(word string) *Position {
	var bestPositions []Position
	maxIntersections := -1

	// Try all possible positions
	for x := 0; x < c.height; x++ {
		for y := 0; y < c.width; y++ {
			for _, dir := range []Direction{Horizontal, Vertical} {
				intersections := c.canBePlaced(word, x, y, dir)
				if intersections < 0 {
					continue
				}

				if intersections > maxIntersections {
					maxIntersections = intersections
					bestPositions = bestPositions[:0]
				}

				if intersections == maxIntersections {
					bestPositions = append(bestPositions, Position{X: x, Y: y, Dir: dir})
				}
			}
		}
	}

	if len(bestPositions) == 0 {
		return nil
	}

	// Return a random position from the best ones
	return &bestPositions[rand.Intn(len(bestPositions))]
}

// GeneratePuzzle generates a crossword puzzle from a list of words
func (c *Crossword) GeneratePuzzle(words []string) bool {
	rand.Seed(time.Now().UnixNano())

	startTime := time.Now()
	maxTime := 1 * time.Minute

	var generate func(pos int) bool
	generate = func(pos int) bool {
		if pos >= len(words) {
			return true
		}

		if time.Since(startTime) > maxTime {
			return false
		}

		word := words[pos]
		if bestPos := c.findBestPosition(word); bestPos != nil {
			// Try placing the word
			c.putWord(word, bestPos.X, bestPos.Y, bestPos.Dir)

			if generate(pos + 1) {
				return true
			}

			// If placing didn't work, remove it and try next position
			c.removeWord(word, bestPos.X, bestPos.Y, bestPos.Dir)
		}

		// Try skipping this word
		return generate(pos + 1)
	}

	return generate(0)
}

// removeWord removes a word from the board
func (c *Crossword) removeWord(word string, x, y int, dir Direction) {
	delete(c.usedWords, word)

	for i := 0; i < len(word); i++ {
		var x1, y1 int
		if dir == Horizontal {
			x1, y1 = x, y+i
			c.hWords[x1][y1] = 0
			if c.vWords[x1][y1] == 0 {
				c.board[x1][y1] = ' '
			}
		} else {
			x1, y1 = x+i, y
			c.vWords[x1][y1] = 0
			if c.hWords[x1][y1] == 0 {
				c.board[x1][y1] = ' '
			}
		}
	}

	// Remove blocking characters if no other words are adjacent
	if dir == Horizontal {
		if c.isValidPosition(x, y-1) && !c.hasAdjacentWords(x, y-1) {
			c.board[x][y-1] = ' '
		}
		if c.isValidPosition(x, y+len(word)) && !c.hasAdjacentWords(x, y+len(word)) {
			c.board[x][y+len(word)] = ' '
		}
	} else {
		if c.isValidPosition(x-1, y) && !c.hasAdjacentWords(x-1, y) {
			c.board[x-1][y] = ' '
		}
		if c.isValidPosition(x+len(word), y) && !c.hasAdjacentWords(x+len(word), y) {
			c.board[x+len(word)][y] = ' '
		}
	}
}

// hasAdjacentWords checks if there are any words adjacent to a position
func (c *Crossword) hasAdjacentWords(x, y int) bool {
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for _, d := range directions {
		newX, newY := x+d[0], y+d[1]
		if c.isValidPosition(newX, newY) && c.board[newX][newY] != ' ' && c.board[newX][newY] != '*' {
			return true
		}
	}
	return false
}

// GetBoard returns the current state of the board
func (c *Crossword) GetBoard() [][]rune {
	return c.board
}

// GetPlacements returns the list of word placements
func (c *Crossword) GetPlacements() []WordPlacement {
	return c.placements
}
