package main

import (
	"fmt"
	"github.com/mohae/deepcopy"
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
	values := strings.Split(string(data), ",")
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

const opAddition = 1
const opMultiplier = 2
const opExit = 99

func newState(opCodePos int64, values []int64) []int64 {
	switch values[opCodePos] {
	case opAddition:
		values[values[opCodePos+3]] = values[values[opCodePos+1]] + values[values[opCodePos+2]]
		return newState(opCodePos+4, values)
	case opMultiplier:
		values[values[opCodePos+3]] = values[values[opCodePos+1]] * values[values[opCodePos+2]]
		return newState(opCodePos+4, values)
	case opExit:
		return values
	default:
		panic(fmt.Sprintf("opcode not found: %d", opCodePos))
	}
}

func bruteForce(values []int64, pos1 int64, pos2 int64) int64 {
	values[1] = pos1
	values[2] = pos2
	v := deepcopy.Copy(values).([]int64)
	v = newState(0, v)
	return v[0]
}
func main() {
	values := readLines()
	toFind := int64(19690720)
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			answer := bruteForce(values, int64(i), int64(j))
			if answer == toFind {
				log.Printf("values %d %d answer: %d", i, j, i*100+j)
			}
		}
	}
}
