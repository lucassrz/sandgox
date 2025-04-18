package main

import (
	"image/color"
	"math/rand"
)

var CellsTypes = map[CellType]CellData{}

type CellType int64

const (
	Air CellType = iota
	Sand
	Water
	Metal
	WaterGenerator
	BlackHole
)

type Cell struct {
	cellType CellType
	color    color.Color
	isActive bool
}

type CellData struct {
	physic  func(x int, y int, g *Game)
	liquid  bool
	density int
}

func initCellsTypes() {
	CellsTypes = map[CellType]CellData{
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

func processCellsPhysic(g *Game) {
	bStart := gridSize
	if gridSize%2 == 0 {
		bStart = bStart - 1
	}
	for yA := 0; yA < gridSize; yA += 1 {
		for xA := 0; xA < gridSize; xA += 2 {
			cellA := g.grid[yA][xA]
			CellsTypes[cellA.cellType].physic(xA, yA, g)
		}

		for xB := bStart; xB > 0; xB -= 2 {
			cellB := g.grid[yA][xB]
			CellsTypes[cellB.cellType].physic(xB, yA, g)
		}
	}

	for yB := bStart; yB > 0; yB -= 2 {
		for xA := 0; xA < gridSize; xA += 2 {
			cellA := g.grid[yB][xA]
			CellsTypes[cellA.cellType].physic(xA, yB, g)
		}
		for xB := bStart; xB > 0; xB -= 2 {
			cellB := g.grid[yB][xB]
			CellsTypes[cellB.cellType].physic(xB, yB, g)
		}
	}
}

func (origin Cell) canSwitchWith(target Cell) bool {
	cellTypeDifferent := target.cellType != origin.cellType
	dataOrigin := CellsTypes[origin.cellType]
	dataTarget := CellsTypes[target.cellType]
	hasOneLiquid := dataTarget.liquid || dataOrigin.liquid
	targetDensityIsInferior := dataTarget.density < dataOrigin.density
	return !origin.isActive && !target.isActive && cellTypeDifferent && (target.cellType == Air || (hasOneLiquid && targetDensityIsInferior))
}

func NewSandCell() Cell {
	colors := []color.Color{
		color.RGBA{255, 255, 0, 255},
		color.RGBA{200, 200, 0, 255},
		color.RGBA{150, 150, 0, 255},
	}
	index := rand.Intn(len(colors))
	if onlyOneColor {
		index = 0
	}
	return Cell{
		cellType: Sand,
		color:    colors[index],
		isActive: true,
	}
}

func SandPhysic(x int, y int, g *Game) {
	cell := g.grid[y][x]
	var actions = make([]func(), 0)
	if y+1 < gridSize {
		if cell.canSwitchWith(g.grid[y+1][x]) {
			actions = append(actions, func() {
				switchPlace(x, y, x, y+1, g)
			})
		}
		if len(actions) == 0 {

			if x-1 >= 0 && cell.canSwitchWith(g.grid[y+1][x-1]) {
				actions = append(actions, func() {
					switchPlace(x, y, x-1, y+1, g)
				})
			}
			if x+1 < gridSize && cell.canSwitchWith(g.grid[y+1][x+1]) {
				actions = append(actions, func() {
					switchPlace(x, y, x+1, y+1, g)
				})
			}
		}

		if len(actions) != 0 {
			actions[rand.Intn(len(actions))]()
		}
	}
}

func switchPlace(Ax int, Ay int, Bx int, By int, g *Game) {
	cellA := g.grid[Ay][Ax]
	cellB := g.grid[By][Bx]
	cellA.isActive = true
	cellB.isActive = true
	g.grid[By][Bx] = cellA
	g.grid[Ay][Ax] = cellB
}

func NewWaterCell() Cell {
	colors := []color.Color{
		color.RGBA{0, 0, 255, 255},
		color.RGBA{0, 0, 200, 255},
		color.RGBA{0, 0, 150, 255},
	}
	index := rand.Intn(len(colors))
	if onlyOneColor {
		index = 0
	}
	return Cell{
		cellType: Water,
		color:    colors[index],
		isActive: true,
	}
}

func WaterPhysic(x int, y int, g *Game) {
	if y+1 < gridSize {
		cell := g.grid[y][x]
		var actions = make([]func(), 0)

		if cell.canSwitchWith(g.grid[y+1][x]) {
			actions = append(actions, func() {
				switchPlace(x, y, x, y+1, g)
			})
		}

		if x+1 < gridSize && cell.canSwitchWith(g.grid[y+1][x+1]) {
			actions = append(actions, func() {
				switchPlace(x, y, x+1, y+1, g)
			})
		}

		if x-1 >= 0 && cell.canSwitchWith(g.grid[y+1][x-1]) {
			actions = append(actions, func() {
				switchPlace(x, y, x-1, y+1, g)
			})
		}

		if len(actions) == 0 {
			if x+1 < gridSize && cell.canSwitchWith(g.grid[y][x+1]) {
				actions = append(actions, func() {
					switchPlace(x, y, x+1, y, g)
				})
			}
			if x-1 >= 0 && cell.canSwitchWith(g.grid[y][x-1]) {
				actions = append(actions, func() {
					switchPlace(x, y, x-1, y, g)
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
		isActive: true,
	}
}

func NewMetalCell() Cell {
	return Cell{
		cellType: Metal,
		color:    color.RGBA{128, 128, 128, 255},
		isActive: true,
	}
}

func NewBlackHoleCell() Cell {
	return Cell{
		cellType: BlackHole,
		color:    color.RGBA{52, 8, 54, 255},
		isActive: true,
	}
}

func NewWaterGeneratorCell() Cell {
	return Cell{
		cellType: WaterGenerator,
		color:    color.RGBA{95, 78, 158, 255},
		isActive: true,
	}
}

func NoPhysic(int, int, *Game) {
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
