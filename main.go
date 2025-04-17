package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ebitenui/ebitenui"
	image2 "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
	"math/rand"
	"os"
)

const (
	screenWidth  = 500
	screenHeight = 500
	cellSize     = 5
	gridSize     = screenWidth / cellSize
)

type Game struct {
	grid             [gridSize][gridSize]Cell
	ui               *ebitenui.UI
	selectedCellType CellType
	brushSize        int
}

type CellType int64

const (
	Air CellType = iota
	Sand
	Water
	Metal
	WaterGenerator
	BlackHole
)

var benchmarkMode = false
var countUpdate = 0

type Cell struct {
	cellType CellType
	color    color.Color
}

type CellData struct {
	physic  func(x int, y int, g *Game)
	liquid  bool
	density int
}

var types = map[CellType]CellData{}

func initTypes() {
	types = map[CellType]CellData{
		Sand: {
			physic:  SandPhysic,
			liquid:  false,
			density: 10,
		},
		Water: {
			physic:  WaterPhysic,
			liquid:  true,
			density: 9,
		},
		Air: {
			physic:  NoPhysic,
			liquid:  false,
			density: 0,
		},
		Metal: {
			physic:  NoPhysic,
			liquid:  false,
			density: 9999,
		},
		BlackHole: {
			physic:  BlackHolePhysic,
			liquid:  false,
			density: 9999,
		},
		WaterGenerator: {
			physic:  WaterGeneratorPhysic,
			liquid:  false,
			density: 9999,
		},
	}
}

func BlackHolePhysic(x int, y int, g *Game) {
	// destroy all cells around black hole

	for offsetY := -1; offsetY <= 1; offsetY++ {
		for offsetX := -1; offsetX <= 1; offsetX++ {
			targetX := x + offsetX
			targetY := y + offsetY
			if targetX >= 0 && targetX < gridSize && targetY >= 0 && targetY < gridSize {
				if g.grid[targetY][targetX].cellType != BlackHole {
					g.grid[targetY][targetX] = NewAirCell()
				}
			}
		}
	}
}

func WaterGeneratorPhysic(x int, y int, g *Game) {
	// generate water all around cell
	for offsetY := -1; offsetY <= 1; offsetY++ {
		for offsetX := -1; offsetX <= 1; offsetX++ {
			targetX := x + offsetX
			targetY := y + offsetY
			if targetX >= 0 && targetX < gridSize && targetY >= 0 && targetY < gridSize {
				if g.grid[targetY][targetX].cellType == Air {
					g.grid[targetY][targetX] = NewWaterCell()
				}
			}
		}
	}

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
			Left:   10,
			Right:  10,
			Top:    10,
			Bottom: 10,
		},
	}
}

func (g *Game) Update() error {

	//var startTime = time.Now()
	g.ui.Update()
	// create a time variable

	// convert to unix time in milliseconds
	for y, row := range g.grid {
		for x, cell := range row {
			types[cell.cellType].physic(x, y, g)
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if x >= 100 {
			cellX := x / cellSize
			cellY := y / cellSize

			for offsetY := -g.brushSize; offsetY <= g.brushSize; offsetY++ {
				for offsetX := -g.brushSize; offsetX <= g.brushSize; offsetX++ {
					targetX := cellX + offsetX
					targetY := cellY + offsetY
					if targetX >= 0 && targetX < gridSize && targetY >= 0 && targetY < gridSize {
						switch g.selectedCellType {
						case Sand:
							g.grid[targetY][targetX] = NewSandCell()
						case Water:
							g.grid[targetY][targetX] = NewWaterCell()
						case Air:
							g.grid[targetY][targetX] = NewAirCell()
						case Metal:
							g.grid[targetY][targetX] = NewMetalCell()
						case BlackHole:
							g.grid[targetY][targetX] = NewBlackHoleCell()
						case WaterGenerator:
							g.grid[targetY][targetX] = NewWaterGeneratorCell()
						}

					}
				}
			}
		}
	}

	//println("Time taken for update: ", time.Since(startTime).Nanoseconds(), "ns")
	return nil
}

func benchmarkCheck() {
	countUpdate++
	if countUpdate >= 100 {
		os.Exit(0)
	}
}

type Rect struct {
	x int
	y int
	w int
	h int
}

func (g *Game) Draw(screen *ebiten.Image) {
	pointsByColor := make(map[color.Color][][]bool)
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			cell := g.grid[y][x]
			if pointsByColor[cell.color] == nil {
				pointsByColor[cell.color] = make([][]bool, gridSize)
				for i := range pointsByColor[cell.color] {
					pointsByColor[cell.color][i] = make([]bool, gridSize)
					for j := range pointsByColor[cell.color][i] {
						pointsByColor[cell.color][i][j] = false
					}
				}
			}
			pointsByColor[cell.color][y][x] = true
		}
	}

	rectanglesByColor := make(map[color.Color][]Rect)
	verticalRectangles(pointsByColor, rectanglesByColor)
	//horizontalRectangles(pointsByColor, rectanglesByColor)

	op := &ebiten.DrawImageOptions{}
	for col, rects := range rectanglesByColor {
		for _, rectangle := range rects {
			rect := getRectImageByWidth(rectangle.w)
			rect.Fill(col)
			op.GeoM.Reset()
			op.GeoM.Translate(float64(rectangle.x*cellSize), float64(rectangle.y*cellSize))
			screen.DrawImage(rect, op)
		}
	}
	g.ui.Draw(screen)
	ebiten.SetVsyncEnabled(false)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
	//println("Time taken for draw: ", time.Since(startTime).Nanoseconds(), "ns")

	if benchmarkMode {
		benchmarkCheck()
	}
}

func verticalRectangles(pointsByColor map[color.Color][][]bool, rectanglesByColor map[color.Color][]Rect) {

	for col, points := range pointsByColor {
		for y, row := range points {
			posX := 0
			width := 0
			for x, point := range row {
				if point {
					if width == 0 {
						posX = x
						width = 1
					} else {
						width++
					}
				} else if width > 0 {
					rectanglesByColor[col] = append(rectanglesByColor[col], Rect{
						x: posX,
						y: y,
						w: width,
						h: 1,
					})
					width = 0
				} else {
					width = 0

				}
			}
			if width > 0 {
				rectanglesByColor[col] = append(rectanglesByColor[col], Rect{
					x: posX,
					y: y,
					w: width,
					h: 1,
				})
			}
		}
	}
}

var cachedRects = make(map[int]*ebiten.Image)

func getRectImageByWidth(width int) *ebiten.Image {
	if rect, ok := cachedRects[width]; ok {
		return rect
	}
	rect := ebiten.NewImage(width*cellSize, cellSize)
	cachedRects[width] = rect
	return rect
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	benchmarkModeUnparsed := flag.Bool("benchmark", false, "benchmark mode")
	flag.Parse()
	benchmarkMode = *benchmarkModeUnparsed
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("sandgox")
	initTypes()
	game := &Game{
		grid:             initGrid(),
		selectedCellType: Sand,
		brushSize:        0,
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

	blackHoleButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Image(res.buttonImage),
		widget.ButtonOpts.Text("Black Hole", res.font, res.textColor),
		widget.ButtonOpts.TextPadding(res.padding),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
			g.selectedCellType = BlackHole
		}),
	)

	waterGeneratorButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Image(res.buttonImage),
		widget.ButtonOpts.Text("Water gen", res.font, res.textColor),
		widget.ButtonOpts.TextPadding(res.padding),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
			g.selectedCellType = WaterGenerator
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
	buttonContainer.AddChild(blackHoleButton)
	buttonContainer.AddChild(waterGeneratorButton)

	brushButtonContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		),
		), widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 0)))

	smallBrushButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Image(res.buttonImage),
		widget.ButtonOpts.Text("1", res.font, res.textColor),
		widget.ButtonOpts.TextPadding(res.padding),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
			g.brushSize = 0
		}),
	)

	mediumBrushButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Image(res.buttonImage),
		widget.ButtonOpts.Text("3", res.font, res.textColor),
		widget.ButtonOpts.TextPadding(res.padding),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
			g.brushSize = 1
		}),
	)

	largeBrushButton := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true})),
		widget.ButtonOpts.Image(res.buttonImage),
		widget.ButtonOpts.Text("7", res.font, res.textColor),
		widget.ButtonOpts.TextPadding(res.padding),
		widget.ButtonOpts.ClickedHandler(func(*widget.ButtonClickedEventArgs) {
			g.brushSize = 3
		}),
	)

	brushButtonContainer.AddChild(smallBrushButton)
	brushButtonContainer.AddChild(mediumBrushButton)
	brushButtonContainer.AddChild(largeBrushButton)
	buttonContainer.AddChild(brushButtonContainer)
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
			if benchmarkMode && y < 10 {
				grid[y][x] = NewWaterCell()
			} else {
				grid[y][x] = NewAirCell()
			}
		}

	}
	return grid
}

func NewSandCell() Cell {
	colors := []color.Color{
		//color.RGBA{255, 255, 0, 255},
		//color.RGBA{200, 200, 0, 255},
		color.RGBA{150, 150, 0, 255},
	}
	return Cell{
		cellType: Sand,
		color:    colors[rand.Intn(len(colors))],
	}
}

func SandPhysic(x int, y int, g *Game) {
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
}

func NewWaterCell() Cell {
	colors := []color.Color{
		//color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 0, 200, 255},
		//color.RGBA{0, 0, 150, 255},
	}
	return Cell{
		cellType: Water,
		color:    colors[rand.Intn(len(colors))],
	}
}

func WaterPhysic(x int, y int, g *Game) {
	if y+1 < gridSize {
		cell := g.grid[y][x]
		var actions = make([]func(), 0)

		if canSwitchCell(cell, g.grid[y+1][x]) {
			actions = append(actions, func() {
				copyOfCell := g.grid[y+1][x]
				g.grid[y+1][x] = cell
				g.grid[y][x] = copyOfCell
			})
		}

		if x+1 < gridSize && canSwitchCell(cell, g.grid[y+1][x+1]) {
			actions = append(actions, func() {
				copyOfCell := g.grid[y+1][x+1]
				g.grid[y+1][x+1] = cell
				g.grid[y][x] = copyOfCell
			})
		}

		if x-1 >= 0 && canSwitchCell(cell, g.grid[y+1][x-1]) {
			actions = append(actions, func() {
				copyOfCell := g.grid[y+1][x-1]
				g.grid[y+1][x-1] = cell
				g.grid[y][x] = copyOfCell
			})
		}

		if len(actions) == 0 {
			if x+1 < gridSize && canSwitchCell(cell, g.grid[y][x+1]) {
				actions = append(actions, func() {
					copyOfCell := g.grid[y][x+1]
					g.grid[y][x+1] = cell
					g.grid[y][x] = copyOfCell
				})
			}
			if x-1 >= 0 && canSwitchCell(cell, g.grid[y][x-1]) {
				actions = append(actions, func() {
					copyOfCell := g.grid[y][x-1]
					g.grid[y][x-1] = cell
					g.grid[y][x] = copyOfCell
				})
			}
		}

		// execute random action

		if len(actions) > 0 {

			randomIndex := rand.Intn(len(actions))
			actions[randomIndex]()
		}

	}
}

func NewAirCell() Cell {
	return Cell{
		cellType: Air,
		color:    color.RGBA{0, 0, 0, 255},
	}
}

func NewMetalCell() Cell {
	return Cell{
		cellType: Metal,
		color:    color.RGBA{128, 128, 128, 255},
	}
}

func NewBlackHoleCell() Cell {
	return Cell{
		cellType: BlackHole,
		color:    color.RGBA{52, 8, 54, 255},
	}
}

func NewWaterGeneratorCell() Cell {
	return Cell{
		cellType: WaterGenerator,
		color:    color.RGBA{95, 78, 158, 255},
	}
}

func NoPhysic(int, int, *Game) {
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
	dataOrigin := types[origin.cellType]
	dataTarget := types[target.cellType]
	hasOneLiquid := dataTarget.liquid || dataOrigin.liquid
	targetDensityIsInferior := dataTarget.density < dataOrigin.density
	return cellTypeDifferent && (target.cellType == Air || (hasOneLiquid && targetDensityIsInferior))
}
