package main

import "math"

// Camera хранит масштаб и положение карты в окне.
type Camera struct {
	Scale   float64
	OffsetX float64
	OffsetY float64
	ZScale  float64
}

func projectPoint(x, y, z int, zScale float64) (float64, float64) {
	projectedX := float64(x-y) * math.Cos(math.Pi/6)
	projectedY := float64(x+y)*math.Sin(math.Pi/6) - float64(z)*zScale
	return projectedX, projectedY
}

func (c Camera) Project(p Point) (float64, float64) {
	px, py := projectPoint(p.X, p.Y, p.Z, c.ZScale)
	return px*c.Scale + c.OffsetX, py*c.Scale + c.OffsetY
}

func (c *Camera) Fit(mapData Map, width, height int) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if mapData.PointCount() == 0 {
		c.Scale = 1
		c.OffsetX = float64(width) / 2
		c.OffsetY = float64(height) / 2
		return
	}

	first := mapData.Rows[0][0]
	minX, minY := projectPoint(first.X, first.Y, first.Z, c.ZScale)
	maxX, maxY := minX, minY
	for y := 0; y < len(mapData.Rows); y++ {
		for x := 0; x < len(mapData.Rows[y]); x++ {
			p := mapData.Rows[y][x]
			px, py := projectPoint(p.X, p.Y, p.Z, c.ZScale)
			if px < minX {
				minX = px
			}
			if px > maxX {
				maxX = px
			}
			if py < minY {
				minY = py
			}
			if py > maxY {
				maxY = py
			}
		}
	}

	mapWidth := maxX - minX
	mapHeight := maxY - minY
	if mapWidth == 0 {
		mapWidth = 1
	}
	if mapHeight == 0 {
		mapHeight = 1
	}

	c.Scale = math.Min(float64(width)/mapWidth, float64(height)/mapHeight) * 0.9
	c.OffsetX = float64(width)/2 - (minX+maxX)/2*c.Scale
	c.OffsetY = float64(height)/2 - (minY+maxY)/2*c.Scale
}
