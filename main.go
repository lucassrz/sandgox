package main

import (
	"flag"
	"fmt"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"image/color"
	"log"
	"os"
)

const (
	screenWidth  = 500
	screenHeight = 500
	menuWidth    = 100
	cellSize     = 5
	gridSize     = screenWidth / cellSize
)

type Game struct {
	grid             [gridSize][gridSize]Cell
	ui               *ebitenui.UI
	selectedCellType CellType
	pixelsToDraw     map[color.Color][][]bool
	brushSize        int
	screenBuffer     *image.RGBA
}

var benchmarkMode = false
var countUpdate = 0
var isChangingBrush = false
var changingBrushTime = 0

type resources struct {
	buttonImage *widget.ButtonImage
	font        text.Face
	textColor   *widget.ButtonTextColor
	sliderImage *widget.SliderTrackImage
	padding     widget.Insets
}

func (g *Game) Update() error {
	g.ui.Update()
	processCellsPhysic(g)
	handleClick(g)
	return nil
}

func benchmarkCheck() {
	countUpdate++
	if countUpdate >= 100 {
		os.Exit(0)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	createScreenBufferImgIfNotExist()
	drawCells(g)
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(screenBufferImg, op)
	drawBrushSize(screen, g)
	g.ui.Draw(screen)
	ebiten.SetVsyncEnabled(false)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
	if benchmarkMode {
		benchmarkCheck()
	}
}

func createScreenBufferImgIfNotExist() {
	if screenBufferImg == nil {
		screenBufferImg = ebiten.NewImage(screenWidth, screenHeight)
		screenBufferImg.Fill(color.RGBA{})
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth + menuWidth, screenHeight
}

func main() {
	initFlags()
	initWindow()
	initCellsTypes()
	game := getGame()
	setupUI(game)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func getGame() *Game {
	return &Game{
		grid:             initGrid(),
		pixelsToDraw:     make(map[color.Color][][]bool),
		selectedCellType: Sand,
		brushSize:        0,
	}
}

func initFlags() {
	benchmarkModeUnparsed := flag.Bool("benchmark", false, "benchmark mode")
	flag.Parse()
	benchmarkMode = *benchmarkModeUnparsed
}

func initWindow() {
	ebiten.SetWindowSize(screenWidth+menuWidth, screenHeight)
	ebiten.SetWindowTitle("sandgox")
}

// Remaining original code (initGrid, cell types, physics) unchanged...
func initGrid() [gridSize][gridSize]Cell {
	grid := [gridSize][gridSize]Cell{}
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			if benchmarkMode && y < 10 {
				grid[y][x] = NewSandCell()
			} else if benchmarkMode && y > 80 && x > 40 && x < 60 {
				grid[y][x] = NewWaterCell()
			} else if benchmarkMode && y == 50 && x > 20 && x < 40 {
				grid[y][x] = NewMetalCell()
			} else if benchmarkMode && y == 50 && x > 60 && x < 95 {
				grid[y][x] = NewBlackHoleCell()
			} else if benchmarkMode && y == 30 && x > 74 && x < 78 {
				grid[y][x] = NewWaterGeneratorCell()
			} else {
				grid[y][x] = NewAirCell()
			}
		}

	}
	return grid
}
