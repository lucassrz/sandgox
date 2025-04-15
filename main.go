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
	deltaTime    = 1.0 / 12.0
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
)

var drawn bool = false
var timeBetweenUpdates = 0
var benchmarkMode bool = false

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

type UpdateCell struct {
	updateFunc func(x int, y int, g *Game)
	PosX       int
	PosY       int
}

func (g *Game) Update() error {

	//var startTime = time.Now()
	g.ui.Update()
	// create a time variable

	// convert to unix time in milliseconds
	if drawn {
		for y, row := range g.grid {
			for x, cell := range row {
				if cell.cellType != Air {
					types[cell.cellType].physic(x, y, g)
				}
			}
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
						}
					}
				}
			}
		}
	}

	if benchmarkMode {
		benchmarkCheck(g)
	}

	//println("Time taken for update: ", time.Since(startTime).Milliseconds(), "ms")
	return nil
}

func benchmarkCheck(game *Game) {
	lineIsFull := true
	for i := 0; i < gridSize; i++ {
		if game.grid[gridSize-1][i].cellType == Air {
			lineIsFull = false
			break
		}
	}

	if lineIsFull {
		os.Exit(0)
	}
}

type ImageToDraw struct {
	image *ebiten.Image
	op    *ebiten.DrawImageOptions
}

func (g *Game) Draw(screen *ebiten.Image) {
	//var startTime = time.Now()
	screen.Fill(color.Black)

	// Pré-créer une image pour les cellules
	rect := ebiten.NewImage(cellSize, cellSize)
	op := &ebiten.DrawImageOptions{}
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			cell := g.grid[y][x]
			if cell.cellType != Air {
				rect.Fill(cell.color) // Remplir avec la couleur de la cellule
				op.GeoM.Reset()       // Réinitialiser les transformations
				op.GeoM.Translate(float64(x*cellSize), float64(y*cellSize))
				screen.DrawImage(rect, op)
			}
		}
	}
	g.ui.Draw(screen)
	ebiten.SetVsyncEnabled(false)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
	drawn = true
	//println("Time taken for draw: ", time.Since(startTime).Milliseconds(), "ms")
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
		color.RGBA{255, 255, 0, 255},
		color.RGBA{200, 200, 0, 255},
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
		color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 0, 200, 255},
		color.RGBA{0, 0, 150, 255},
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

func NoPhysic(x int, y int, g *Game) {
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

func hasADifferentNeighbor(x int, y int, cell Cell, g *Game) bool {
	cellType := cell.cellType
	rightIsOk := x+1 < gridSize
	grid := g.grid
	if rightIsOk && grid[y][x+1].cellType != cellType {
		return true
	}
	leftIsOk := x-1 >= 0
	if leftIsOk && grid[y][x-1].cellType != cellType {
		return true
	}
	bottomIsOk := y+1 < gridSize
	if bottomIsOk && grid[y+1][x].cellType != cellType {
		return true
	}
	topIsOk := y-1 >= 0
	if topIsOk && grid[y-1][x].cellType != cellType {
		return true
	}
	if rightIsOk && bottomIsOk && grid[y+1][x+1].cellType != cellType {
		return true
	}
	if leftIsOk && bottomIsOk && grid[y+1][x-1].cellType != cellType {
		return true
	}
	if rightIsOk && topIsOk && grid[y-1][x+1].cellType != cellType {
		return true
	}
	if leftIsOk && topIsOk && grid[y-1][x-1].cellType != cellType {
		return true
	}
	return false
}
