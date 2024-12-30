// File: utils/render.go
package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// RenderConfig holds configuration for rendering the crossword
type RenderConfig struct {
	CellSize        int     // Size of each cell in pixels
	BorderSize      int     // Size of border lines
	FontSize        float64 // Font size for letters
	BackgroundColor color.Color
	GridLineColor   color.Color
	BlockColor      color.Color
	LetterColor     color.Color
}

// DefaultConfig returns a default rendering configuration
func DefaultConfig() RenderConfig {
	return RenderConfig{
		CellSize:        40,
		BorderSize:      2,
		FontSize:        24,
		BackgroundColor: color.White,
		GridLineColor:   color.Black,
		BlockColor:      color.Black,
		LetterColor:     color.Black,
	}
}

// RenderPuzzleToPNG creates a PNG image of the crossword puzzle
func RenderPuzzleToPNG(puzzle *Crossword, filename string, config RenderConfig) error {
	board := puzzle.GetBoard()
	height := len(board)
	width := len(board[0])

	// Calculate image dimensions
	imgWidth := width*config.CellSize + config.BorderSize
	imgHeight := height*config.CellSize + config.BorderSize

	// Create new image
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{config.BackgroundColor}, image.Point{}, draw.Src)

	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return err
	}

	// Create font context
	fontContext := freetype.NewContext()
	fontContext.SetDPI(72)
	fontContext.SetFont(font)
	fontContext.SetFontSize(config.FontSize)
	fontContext.SetClip(img.Bounds())
	fontContext.SetDst(img)
	fontContext.SetSrc(image.NewUniform(config.LetterColor))

	// Draw grid and fill cells
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cell := board[y][x]
			cellX := x * config.CellSize
			cellY := y * config.CellSize

			// Draw cell border
			drawRect(img, cellX, cellY, config.CellSize, config.CellSize, config.GridLineColor)

			// Fill black squares for blocked cells
			if cell == '*' {
				fillRect(img,
					cellX+config.BorderSize,
					cellY+config.BorderSize,
					config.CellSize-2*config.BorderSize,
					config.CellSize-2*config.BorderSize,
					config.BlockColor)
			} else if cell != ' ' {
				// Draw letter
				letter := strings.ToUpper(string(cell))

				// Calculate text position (centered in cell)
				textWidth := config.FontSize * 0.6 // Approximate width of character
				textX := float64(cellX) + (float64(config.CellSize)-textWidth)/2
				textY := float64(cellY) + float64(config.CellSize)*0.7 // Adjust for baseline

				fontContext.SetDst(img)
				fontContext.SetClip(image.Rect(cellX, cellY, cellX+config.CellSize, cellY+config.CellSize))
				_, err := fontContext.DrawString(letter, freetype.Pt(int(textX), int(textY)))
				if err != nil {
					return err
				}
			}
		}
	}

	// Add numbers for word starts
	placements := puzzle.GetPlacements()
	wordNumbers := make(map[string]int)
	currentNumber := 1

	for _, placement := range placements {
		key := fmt.Sprintf("%d,%d", placement.X, placement.Y)
		if _, exists := wordNumbers[key]; !exists {
			wordNumbers[key] = currentNumber

			// Draw number
			numberStr := fmt.Sprintf("%d", currentNumber)
			fontContext.SetFontSize(config.FontSize * 0.4)
			fontContext.DrawString(numberStr,
				freetype.Pt(
					placement.Y*config.CellSize+config.BorderSize+2,
					placement.X*config.CellSize+config.BorderSize+10))

			currentNumber++
		}
	}

	// Save to file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// Helper function to draw a rectangle outline
func drawRect(img *image.RGBA, x, y, w, h int, c color.Color) {
	// Top
	drawHLine(img, x, y, w, c)
	// Bottom
	drawHLine(img, x, y+h-1, w, c)
	// Left
	drawVLine(img, x, y, h, c)
	// Right
	drawVLine(img, x+w-1, y, h, c)
}

// Helper function to fill a rectangle
func fillRect(img *image.RGBA, x, y, w, h int, c color.Color) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, c)
		}
	}
}

// Helper function to draw horizontal line
func drawHLine(img *image.RGBA, x, y, w int, c color.Color) {
	for i := 0; i < w; i++ {
		img.Set(x+i, y, c)
	}
}

// Helper function to draw vertical line
func drawVLine(img *image.RGBA, x, y, h int, c color.Color) {
	for i := 0; i < h; i++ {
		img.Set(x, y+i, c)
	}
}
