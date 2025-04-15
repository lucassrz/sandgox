package main

import (
	"bytes"
	"fmt"
	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
	"math/rand"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 500
	screenHeight = 500
	cellSize     = 5
	gridSize     = screenWidth / cellSize
	deltaTime    = 1.0 / 60.0
)

type Game struct {
	grid             [gridSize][gridSize]Cell
	ui               *ebitenui.UI
	selectedCellType CellType
}

type CellType int64

const (
	Air CellType = iota
	Sand
	Water
	Metal
)

type Cell struct {
	physic   func(x int, y int, g *Game)
	cellType CellType
	color    color.Color
	liquid   bool
	density  int
}

type resources struct {
	buttonImage *widget.ButtonImage
	font        text.Face
	textColor   *widget.ButtonTextColor
	padding     widget.Insets
}

func newResources() *resources {
	idle := image2.NewNineSliceColor(color.NRGBA{R: 0x33, G: 0x33, B: 0x33, A: 0xff})
	hover := image2.NewNineSliceColor(color.NRGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff})
	pressed := image2.NewNineSliceColor(color.NRGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff})

	font, _ := loadFont(10)
	return &resources{
		buttonImage: &widget.ButtonImage{
			Idle:    idle,
			Hover:   hover,
			Pressed: pressed,
		},
		textColor: &widget.ButtonTextColor{
			Idle: color.NRGBA{R: 0xdf, G: 0xf4, B: 0xff, A: 0xff},
		},
		font: font,
		padding: widget.Insets{
			Left:   30,
			Right:  30,
			Top:    10,
			Bottom: 10,
		},
	}
}

func (g *Game) Update() error {
	g.ui.Update()
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			g.grid[y][x].physic(x, y, g)
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x >= 100 {
			cellX := x / cellSize
			cellY := y / cellSize
			if cellX >= 0 && cellX < gridSize && cellY >= 0 && cellY < gridSize {
				switch g.selectedCellType {
				case Sand:
					g.grid[cellY][cellX] = NewSandCell()
				case Water:
					g.grid[cellY][cellX] = NewWaterCell()
				case Air:
					g.grid[cellY][cellX] = NewAirCell()
				case Metal:
					g.grid[cellY][cellX] = NewMetalCell()
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			cell := g.grid[y][x]
			if cell.cellType != Air {
				rect := ebiten.NewImage(cellSize, cellSize)
				rect.Fill(cell.color)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*cellSize), float64(y*cellSize))
				screen.DrawImage(rect, op)
			}
		}
	}

	g.ui.Draw(screen)
	ebiten.SetVsyncEnabled(false)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pixel Sand Simulation")
	game := &Game{
		grid:             initGrid(),
		selectedCellType: Sand,
	}
	game.setupUI()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) setupUI() {
	res := newResources()
	sandButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Image(res.buttonImage),
		widget.ButtonOpts.Text("Sand", res.font, res.textColor),
		widget.ButtonOpts.TextPadding(res.padding),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
			g.selectedCellType = Sand
		}),
	)

	waterButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Image(res.buttonImage),
		widget.ButtonOpts.Text("Water", res.font, res.textColor),
		widget.ButtonOpts.TextPadding(res.padding),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
			g.selectedCellType = Water
		}),
	)

	airButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Image(res.buttonImage),
		widget.ButtonOpts.Text("Air", res.font, res.textColor),
		widget.ButtonOpts.TextPadding(res.padding),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
			g.selectedCellType = Air
		}),
	)

	metalButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Image(res.buttonImage),
		widget.ButtonOpts.Text("Metal", res.font, res.textColor),
		widget.ButtonOpts.TextPadding(res.padding),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
			g.selectedCellType = Metal
		}),
	)

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		),
		), widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 0)))

	buttonContainer.AddChild(sandButton)
	buttonContainer.AddChild(waterButton)
	buttonContainer.AddChild(airButton)
	buttonContainer.AddChild(metalButton)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	rootContainer.AddChild(buttonContainer)
	g.ui = &ebitenui.UI{
		Container: rootContainer,
	}
}

// Remaining original code (initGrid, cell types, physics) unchanged...
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
			fmt.Println("Sand physic" + fmt.Sprintf("x: %d, y: %d", x, y))
			if y+1 < gridSize {
				cell := g.grid[y][x]
				if canSwitchCell(cell, g.grid[y+1][x]) {
					copyOfCell := g.grid[y+1][x]
					g.grid[y+1][x] = cell
					g.grid[y][x] = copyOfCell
				} else if x-1 >= 0 && canSwitchCell(cell, g.grid[y+1][x-1]) {
					copyOfCell := g.grid[y+1][x-1]
					g.grid[y+1][x-1] = cell
					g.grid[y][x] = copyOfCell
				} else if x+1 < gridSize && canSwitchCell(cell, g.grid[y+1][x+1]) {
					copyOfCell := g.grid[y+1][x+1]
					g.grid[y+1][x+1] = cell
					g.grid[y][x] = copyOfCell
				}
			}
		},
		cellType: Sand,
		liquid:   false,
		density:  10,
		color:    colors[rand.Intn(len(colors))],
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
			if y+1 < gridSize && g.grid[y+1][x].cellType == Air {
				g.grid[y+1][x] = g.grid[y][x]
				g.grid[y][x] = NewAirCell()
			}
		},
		cellType: Water,
		liquid:   true,
		density:  9,
		color:    colors[rand.Intn(len(colors))],
	}
}

func NewAirCell() Cell {
	return Cell{
		physic:   func(x int, y int, g *Game) {},
		cellType: Air,
		liquid:   false,
		density:  0,
		color:    color.RGBA{0, 0, 0, 255},
	}
}

func NewMetalCell() Cell {
	return Cell{
		physic: func(x int, y int, g *Game) {
		},
		cellType: Metal,
		liquid:   false,
		density:  9999,
		color:    color.RGBA{128, 128, 128, 255},
	}
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func canSwitchCell(origin Cell, target Cell) bool {
	cellTypeDifferent := target.cellType != origin.cellType
	hasOneLiquid := target.liquid || origin.liquid
	targetDensityIsInferior := target.density < origin.density
	return cellTypeDifferent && (target.cellType == Air || (hasOneLiquid && targetDensityIsInferior))
}
