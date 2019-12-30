package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func readLines() []int64 {
	f, err := os.Open("../data")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic("err")
	}
	values := strings.Split(string(data), "")
	intVals := make([]int64, 0)
	for _, value := range values {
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(err)
		}
		intVals = append(intVals, intVal)
	}
	return intVals
}

type Layer struct {
	Lines [][]int64
}

func generateLayers(lines []int64, width int, height int) []Layer {
	layers := make([]Layer, 0)
	for i := 0; i < len(lines); i += width * height {
		l := Layer{Lines: make([][]int64, 0)}
		for k := 0; k < height; k++ {
			currentStart := i + k*width
			line := lines[currentStart : currentStart+width]

			l.Lines = append(l.Lines, line)
		}
		layers = append(layers, l)
	}
	return layers
}

const transparent = 2
const black = 0

func main() {
	width := 25
	height := 6
	values := readLines()
	layers := generateLayers(values, width, height)
	finalLayer := Layer{
		Lines: make([][]int64, height),
	}

	for _, layer := range layers {
		for i := 0; i < len(layer.Lines); i++ {
			if len(finalLayer.Lines[i]) == 0 {
				finalLayer.Lines[i] = make([]int64, width)
				for idx, _ := range finalLayer.Lines[i] {
					finalLayer.Lines[i][idx] = transparent
				}
			}
			for j := 0; j < len(layer.Lines[i]); j++ {
				if finalLayer.Lines[i][j] != transparent {
					continue
				}
				if layer.Lines[i][j] != transparent {
					finalLayer.Lines[i][j] = layer.Lines[i][j]
				}
			}
		}
	}

	// ugly formatting code
	for _, line := range finalLayer.Lines {
		formatted := ""
		for _, digit := range line {
			if digit == black {
				formatted += "  "
			} else {
				formatted += "x "
			}
		}
		log.Printf("%s", formatted)
	}
}
