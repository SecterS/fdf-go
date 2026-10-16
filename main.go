package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Использование: go run . <файл.fdf>")
		os.Exit(2)
	}

	mapData, err := readMap(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка карты:", err)
		os.Exit(1)
	}

	fmt.Println("Карта прочитана, точек:", mapData.PointCount())
}
