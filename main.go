package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 500
	screenHeight = 500
	cellSize     = 5
	gridSize     = screenWidth / cellSize
	deltaTime    = 1.0 / 60.0
	lastUpdate   = 0.0
)

type Game struct {
	grid [gridSize][gridSize]Cell
}
type CellType int64

const (
	Air CellType = iota
	Sand
	Water
	Metal
)

type Cell struct {
	physic   func(x int, y int, g *Game) // Physic function
	cellType CellType
	color    color.Color
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cellX := x / cellSize
		cellY := y / cellSize
		if cellX >= 0 && cellX < gridSize && cellY >= 0 && cellY < gridSize {
			g.grid[cellY][cellX] = NewSandCell()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	ebiten.SetVsyncEnabled(false)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			g.grid[y][x].physic(x, y, g)
		}
	}
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			cell := g.grid[y][x]
			if cell.cellType != Air {
				rect := ebiten.NewImage(cellSize, cellSize)
				rect.Fill(cell.color) // rouge
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
	if err := ebiten.RunGame(&Game{
		grid: initGrid(),
	}); err != nil {
		log.Fatal(err)
	}
}

func initGrid() [gridSize][gridSize]Cell {
	grid := [gridSize][gridSize]Cell{}
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			grid[y][x] = NewAirCell()
		}
	}
	return grid
}

func NewSandCell() Cell {
	colors := []color.Color{
		color.RGBA{255, 255, 0, 255},
		color.RGBA{200, 200, 0, 255},
		color.RGBA{150, 150, 0, 255},
	}
	return Cell{
		physic: func(x int, y int, g *Game) {
			if y+1 < gridSize {
				if g.grid[y+1][x].cellType == Air {
					g.grid[y+1][x] = g.grid[y][x]
					g.grid[y][x] = NewAirCell()
				}
				if x-1 >= 0 && g.grid[y+1][x-1].cellType == Air {
					g.grid[y+1][x-1] = g.grid[y][x]
					g.grid[y][x] = NewAirCell()
				}
				if x+1 < gridSize && g.grid[y+1][x+1].cellType == Air {
					g.grid[y+1][x+1] = g.grid[y][x]
					g.grid[y][x] = NewAirCell()
				}
			}

		},
		cellType: Sand,
		color:    colors[rand.Intn(len(colors))], // yellow
	}
}

func NewWaterCell() Cell {
	colors := []color.Color{
		color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 0, 200, 255},
		color.RGBA{0, 0, 150, 255},
	}
	return Cell{
		physic: func(x int, y int, g *Game) {
			if y+1 < gridSize {
				if g.grid[y+1][x].cellType == Air {
					g.grid[y+1][x] = g.grid[y][x]
					g.grid[y][x] = NewAirCell()
				}
			}

		},
		color:    colors[rand.Intn(len(colors))], // yellow
		cellType: Water,
	}
}

func NewAirCell() Cell {
	colors := []color.Color{
		color.RGBA{0, 0, 0, 255},
	}
	return Cell{
		physic: func(x int, y int, g *Game) {
			// Do nothing for air cells
		},
		cellType: Air,
		color:    colors[rand.Intn(len(colors))], // yellow
	}
}
func NewMetalCell() Cell {
	colors := []color.Color{
		color.RGBA{50, 50, 50, 255},
		color.RGBA{60, 60, 60, 255},
		color.RGBA{70, 70, 70, 255},
	}
	return Cell{
		physic: func(x int, y int, g *Game) {
			// Do nothing for air cells
		},
		cellType: Metal,
		color:    colors[rand.Intn(len(colors))], // yellow
	}
}
