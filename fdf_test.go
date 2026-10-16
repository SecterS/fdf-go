package main

import (
	"errors"
	"math"
	"os"
	"testing"
)

func TestReadElem(t *testing.T) {
	mapData, err := readMap("maps/elem.fdf")
	if err != nil {
		t.Fatal(err)
	}
	if len(mapData.Rows) != 5 || len(mapData.Rows[0]) != 5 {
		t.Fatalf("ожидалась карта 5x5, получено %dx%d", len(mapData.Rows), len(mapData.Rows[0]))
	}
	center := mapData.Rows[2][2]
	if center.X != 2 || center.Y != 2 || center.Z != 5 {
		t.Fatalf("неверная центральная точка: %+v", center)
	}
}

func checkFormatError(t *testing.T, path string, line, column int) {
	t.Helper()
	_, err := readMap(path)
	var formatErr *MapFormatError
	if !errors.As(err, &formatErr) {
		t.Fatalf("ожидалась MapFormatError, получено %v", err)
	}
	if formatErr.Line != line || formatErr.Column != column {
		t.Fatalf("ожидалась строка %d, колонка %d; получено: %v", line, column, err)
	}
}

func TestEmptyFile(t *testing.T) {
	checkFormatError(t, "maps/bad/empty.fdf", 1, 1)
}

func TestBadNumber(t *testing.T) {
	checkFormatError(t, "maps/bad/number.fdf", 2, 3)
}

func TestUnevenRows(t *testing.T) {
	checkFormatError(t, "maps/bad/uneven.fdf", 2, 3)
}

func TestBadColor(t *testing.T) {
	checkFormatError(t, "maps/bad/color.fdf", 2, 2)
}

func TestBlankMiddleRow(t *testing.T) {
	checkFormatError(t, "maps/bad/blank.fdf", 2, 1)
}

func TestParseValue(t *testing.T) {
	tests := []struct {
		word  string
		z     int
		color string
	}{
		{"7", 7, ""},
		{"-3", -3, ""},
		{"2,0xFF8800", 2, "#FF8800"},
		{"1,0xff", 1, "#0000FF"},
	}
	for _, test := range tests {
		t.Run(test.word, func(t *testing.T) {
			point, err := parseValue(test.word, 0, 0)
			if err != nil {
				t.Fatal(err)
			}
			if point.Z != test.z || point.Color != test.color {
				t.Fatalf("для %q получено %+v", test.word, point)
			}
		})
	}
}

func TestMissingFileIsNotFormatError(t *testing.T) {
	_, err := readMap("maps/does-not-exist.fdf")
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ожидалась ошибка отсутствия файла, получено %v", err)
	}
	var formatErr *MapFormatError
	if errors.As(err, &formatErr) {
		t.Fatalf("ошибка файла не должна быть MapFormatError: %v", err)
	}
}

func TestProjectOrigin(t *testing.T) {
	camera := Camera{Scale: 1, ZScale: 1}
	x, y := camera.Project(Point{X: 0, Y: 0, Z: 0})
	if x != 0 || y != 0 {
		t.Fatalf("ожидалось (0, 0), получено (%v, %v)", x, y)
	}
}

func TestProjectHeight(t *testing.T) {
	camera := Camera{Scale: 1, ZScale: 1}
	_, y := camera.Project(Point{X: 0, Y: 0, Z: 1})
	if y >= 0 {
		t.Fatalf("высокая точка должна быть выше нуля, получено y=%v", y)
	}
	camera.ZScale = 0
	_, y = camera.Project(Point{X: 0, Y: 0, Z: 1})
	if y != 0 {
		t.Fatalf("при ZScale=0 высота не должна влиять на y, получено %v", y)
	}
}

func TestCameraFit(t *testing.T) {
	mapData, err := readMap("maps/elem.fdf")
	if err != nil {
		t.Fatal(err)
	}
	camera := Camera{ZScale: 1}
	camera.Fit(mapData, 800, 600)
	for _, point := range mapData.Points() {
		x, y := camera.Project(point)
		if x < 0 || x > 800 || y < 0 || y > 600 {
			t.Fatalf("точка %+v за пределами окна: (%v, %v)", point, x, y)
		}
	}
}

func TestCameraFitOnePoint(t *testing.T) {
	mapData := Map{Rows: [][]Point{{{X: 0, Y: 0, Z: 4}}}}
	camera := Camera{ZScale: 1}
	camera.Fit(mapData, 800, 600)
	x, y := camera.Project(mapData.Rows[0][0])
	if math.IsNaN(x) || math.IsNaN(y) || x != 400 || y != 300 {
		t.Fatalf("одна точка должна быть в центре, получено (%v, %v)", x, y)
	}
}

func TestPointCountAndPoints(t *testing.T) {
	mapData := Map{Rows: [][]Point{
		{{X: 0, Y: 0}, {X: 1, Y: 0}},
		{{X: 0, Y: 1}, {X: 1, Y: 1}},
	}}
	if mapData.PointCount() != 4 || len(mapData.Points()) != 4 {
		t.Fatalf("ожидалось 4 точки, получено %d", mapData.PointCount())
	}
	if mapData.Points()[2].X != 0 || mapData.Points()[2].Y != 1 {
		t.Fatalf("неверный порядок обхода: %+v", mapData.Points())
	}
}

func TestNeighbors(t *testing.T) {
	mapData := Map{Rows: [][]Point{
		{{X: 0, Y: 0}, {X: 1, Y: 0}},
		{{X: 0, Y: 1}, {X: 1, Y: 1}},
	}}
	neighbors := mapData.Neighbors(mapData.Rows[0][0])
	if len(neighbors) != 2 || neighbors[0].X != 1 || neighbors[0].Y != 0 || neighbors[1].X != 0 || neighbors[1].Y != 1 {
		t.Fatalf("верхний левый угол должен иметь соседа справа и снизу: %+v", neighbors)
	}
	if len(mapData.Neighbors(mapData.Rows[1][1])) != 0 {
		t.Fatal("нижний правый угол не должен иметь соседей справа или снизу")
	}
}

func TestTrailingEmptyLines(t *testing.T) {
	path := t.TempDir() + "/end.fdf"
	if err := os.WriteFile(path, []byte("1 2\n3 4\n\n  \n"), 0600); err != nil {
		t.Fatal(err)
	}
	mapData, err := readMap(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(mapData.Rows) != 2 || mapData.PointCount() != 4 {
		t.Fatalf("пустые строки в конце не должны добавлять точки: %+v", mapData)
	}
}

func TestMalformedColorValues(t *testing.T) {
	for _, word := range []string{"1,0x", "1,0x1234567", "1,0xGG", "1,0xff,0x00"} {
		t.Run(word, func(t *testing.T) {
			_, err := parseValue(word, 2, 1)
			var formatErr *MapFormatError
			if !errors.As(err, &formatErr) || formatErr.Line != 2 || formatErr.Column != 3 {
				t.Fatalf("ожидалась ошибка цвета в строке 2, колонке 3; получено %v", err)
			}
		})
	}
}
