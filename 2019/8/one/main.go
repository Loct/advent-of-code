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
	for i := 0; i < len(lines); i+=width * height {
		l := Layer{Lines:make([][]int64, 0)}
		for k := 0; k < height; k++ {
			currentStart := i + k * width
			line := lines[currentStart: currentStart + width]

			l.Lines = append(l.Lines, line)
		}
		layers = append(layers, l)
	}
	return layers
}

func main() {
	width := 25
	height := 6
	values := readLines()
	layers := generateLayers(values, width, height)
	leastZeros := 100000000
	leastOnes := 0
	leastTwos := 0
	for _, layer := range layers {
		zeros := 0
		ones := 0
		twos := 0
		for _, line := range layer.Lines {
			for _, digit := range line {
				if digit == 0 {
					zeros++
				}
				if digit == 1 {
					ones++
				}
				if digit == 2 {
					twos++
				}
			}
		}
		if leastZeros > zeros {
			leastZeros = zeros
			leastOnes = ones
			leastTwos = twos
		}
	}

	log.Printf("%d %d %d", leastZeros, leastOnes, leastTwos)

}
