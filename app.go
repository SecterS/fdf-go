package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	mapData Map
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	camera := Camera{ZScale: 1}
	bounds := screen.Bounds()
	camera.Fit(g.mapData, bounds.Dx(), bounds.Dy())

	for y := 0; y < len(g.mapData.Rows); y++ {
		for x := 0; x < len(g.mapData.Rows[y]); x++ {
			point := g.mapData.Rows[y][x]
			startX, startY := camera.Project(point)
			for _, neighbor := range g.mapData.Neighbors(point) {
				endX, endY := camera.Project(neighbor)
				vector.StrokeLine(screen, float32(startX), float32(startY), float32(endX), float32(endY), 2, color.White, true)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth < 1 {
		outsideWidth = 1
	}
	if outsideHeight < 1 {
		outsideHeight = 1
	}
	return outsideWidth, outsideHeight
}

func runApp(mapData Map) error {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("FdF")
	return ebiten.RunGame(&Game{mapData: mapData})
}
