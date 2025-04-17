package main

import "github.com/hajimehoshi/ebiten/v2"

func handleClick(g *Game) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		if x <= screenWidth {
			cellX := x / cellSize
			cellY := y / cellSize
			cellConstructor := getCellConstructor(g)
			for offsetY := -g.brushSize; offsetY <= g.brushSize; offsetY++ {
				for offsetX := -g.brushSize; offsetX <= g.brushSize; offsetX++ {
					targetX := cellX + offsetX
					targetY := cellY + offsetY
					if targetX >= 0 && targetX < gridSize && targetY >= 0 && targetY < gridSize {
						if g.selectedCellType == Air || g.grid[targetY][targetX].cellType == Air {
							g.grid[targetY][targetX] = cellConstructor()
						}
					}
				}
			}
		}
	}
}

func getCellConstructor(g *Game) func() Cell {
	cellConstructor := func() Cell {
		return Cell{}
	}
	switch g.selectedCellType {
	case Sand:
		cellConstructor = NewSandCell
	case Water:
		cellConstructor = NewWaterCell
	case Air:
		cellConstructor = NewAirCell
	case Metal:
		cellConstructor = NewMetalCell
	case BlackHole:
		cellConstructor = NewBlackHoleCell
	case WaterGenerator:
		cellConstructor = NewWaterGeneratorCell
	}
	return cellConstructor
}
