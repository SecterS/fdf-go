package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MapFormatError хранит место и причину ошибки в карте.
type MapFormatError struct {
	Line, Column int
	Reason       string
}

func (e *MapFormatError) Error() string {
	return fmt.Sprintf("строка %d, колонка %d: %s", e.Line, e.Column, e.Reason)
}

func parseValue(word string, x, y int) (Point, error) {
	parts := strings.Split(word, ",")
	if len(parts) > 2 {
		return Point{}, &MapFormatError{y + 1, x + 1, "некорректный цвет: " + word}
	}

	z, err := strconv.Atoi(parts[0])
	if err != nil {
		return Point{}, &MapFormatError{y + 1, x + 1, "не целое число: " + word}
	}

	color := ""
	if len(parts) == 2 {
		value := parts[1]
		if len(value) < 3 || len(value) > 8 || strings.ToLower(value[:2]) != "0x" {
			return Point{}, &MapFormatError{y + 1, x + 1, "некорректный цвет: " + word}
		}
		number, err := strconv.ParseUint(value[2:], 16, 24)
		if err != nil {
			return Point{}, &MapFormatError{y + 1, x + 1, "некорректный цвет: " + word}
		}
		color = fmt.Sprintf("#%06X", number)
	}
	return Point{X: x, Y: y, Z: z, Color: color}, nil
}

func readMap(path string) (Map, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Map{}, err
	}

	text := strings.TrimRight(string(data), " \t\r\n")
	if text == "" {
		return Map{}, &MapFormatError{1, 1, "пустой файл"}
	}

	lines := strings.Split(text, "\n")
	mapData := Map{Rows: [][]Point{}}
	width := 0
	for y := 0; y < len(lines); y++ {
		words := strings.Fields(lines[y])
		if len(words) == 0 {
			return Map{}, &MapFormatError{y + 1, 1, "пустая строка"}
		}
		if y == 0 {
			width = len(words)
		} else if len(words) != width {
			column := len(words) + 1
			if len(words) > width {
				column = width + 1
			}
			return Map{}, &MapFormatError{y + 1, column, "строки разной длины"}
		}

		row := []Point{}
		for x := 0; x < len(words); x++ {
			point, err := parseValue(words[x], x, y)
			if err != nil {
				return Map{}, err
			}
			row = append(row, point)
		}
		mapData.Rows = append(mapData.Rows, row)
	}
	return mapData, nil
}
