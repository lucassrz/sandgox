package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 500
	screenHeight = 500
	cellSize     = 5
	gridSize     = screenWidth / cellSize
)

type Game struct {
	grid [gridSize][gridSize]bool
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cellX := x / cellSize
		cellY := y / cellSize
		if cellX >= 0 && cellX < gridSize && cellY >= 0 && cellY < gridSize {
			g.grid[cellY][cellX] = true
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			if g.grid[y][x] {
				rect := ebiten.NewImage(cellSize, cellSize)
				rect.Fill(color.RGBA{255, 0, 0, 255}) // rouge
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*cellSize), float64(y*cellSize))
				screen.DrawImage(rect, op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Grille Pixel Simulation")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
