package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"time"
)

var (
	screenBufferImg *ebiten.Image
)

func groupRectanglesHorizontallyByColor(pointsByColor map[color.Color][][]bool) map[color.Color][]Rect {
	rectanglesByColor := make(map[color.Color][]Rect)
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
	return rectanglesByColor
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

func groupUpdatedCellsByColor(g *Game) map[color.Color][][]bool {
	pointsByColor := make(map[color.Color][][]bool)
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			if g.grid[y][x].isActive {
				if pointsByColor[g.grid[y][x].color] == nil {
					pointsByColor[g.grid[y][x].color] = make([][]bool, gridSize)
					for i := range pointsByColor[g.grid[y][x].color] {
						pointsByColor[g.grid[y][x].color][i] = make([]bool, gridSize)
						for j := range pointsByColor[g.grid[y][x].color][i] {
							pointsByColor[g.grid[y][x].color][i][j] = false
						}
					}
				}
				g.grid[y][x].isActive = false
				pointsByColor[g.grid[y][x].color][y][x] = true
			}
		}
	}
	return pointsByColor
}

func drawBrushSize(screen *ebiten.Image, g *Game) {
	if isChangingBrush {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64((50-g.brushSize)*cellSize), float64((50-g.brushSize)*cellSize))
		rect := ebiten.NewImage((g.brushSize*2+1)*cellSize, (g.brushSize*2+1)*cellSize)
		rect.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})
		screen.DrawImage(rect, op)
		if time.Now().Second() != changingBrushTime {
			isChangingBrush = false
		}
	}
}

func drawCells(g *Game) {
	updatedCellsByColor := groupUpdatedCellsByColor(g)
	rectanglesByColor := groupRectanglesHorizontallyByColor(updatedCellsByColor)
	drawRectangles(rectanglesByColor)
}

func drawRectangles(rectanglesByColor map[color.Color][]Rect) {
	op := &ebiten.DrawImageOptions{}
	for col, rects := range rectanglesByColor {
		for _, rectangle := range rects {
			rect := getRectImageByWidth(rectangle.w)
			rect.Fill(col)
			op.GeoM.Reset()
			op.GeoM.Translate(float64(rectangle.x*cellSize), float64(rectangle.y*cellSize))
			screenBufferImg.DrawImage(rect, op)
		}
	}
}

type Rect struct {
	x int
	y int
	w int
	h int
}
