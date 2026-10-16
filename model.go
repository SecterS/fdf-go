package main

// Point — одна точка карты.
type Point struct {
	X, Y, Z int
	Color   string
}

// Map хранит точки по строкам, как в файле.
type Map struct {
	Rows [][]Point
}

func (m Map) PointCount() int {
	count := 0
	for y := 0; y < len(m.Rows); y++ {
		count += len(m.Rows[y])
	}
	return count
}

func (m Map) Points() []Point {
	points := []Point{}
	for y := 0; y < len(m.Rows); y++ {
		for x := 0; x < len(m.Rows[y]); x++ {
			points = append(points, m.Rows[y][x])
		}
	}
	return points
}

// Neighbors возвращает только соседей справа и снизу.
func (m Map) Neighbors(p Point) []Point {
	if p.Y < 0 || p.Y >= len(m.Rows) || p.X < 0 || p.X >= len(m.Rows[p.Y]) {
		return nil
	}

	neighbors := []Point{}
	if p.X+1 < len(m.Rows[p.Y]) {
		neighbors = append(neighbors, m.Rows[p.Y][p.X+1])
	}
	if p.Y+1 < len(m.Rows) && p.X < len(m.Rows[p.Y+1]) {
		neighbors = append(neighbors, m.Rows[p.Y+1][p.X])
	}
	return neighbors
}
