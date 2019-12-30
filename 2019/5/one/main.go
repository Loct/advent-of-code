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

type operator int64

const opAddition operator = 1
const opMultiplier operator = 2
const opInput operator = 3
const opOutput operator = 4
const opJumpTrue operator = 5
const opJumpFalse operator = 6
const lessThan operator = 7
const equals operator = 8
const opExit operator = 99

type mode int64

const immediateMode mode = 1
const positionMode mode = 0

func getOpcode(input int64) []int64 {
	digits := strconv.FormatInt(input, 10)
	toAppend := 5 - len(digits)
	for i := 0; i < toAppend; i++ {
		digits = "0" + digits
	}
	items := make([]int64, 4)
	var err error
	items[0], err = strconv.ParseInt(string(digits[0]), 10, 64)
	if err != nil {
		panic(err)
	}
	items[1], err = strconv.ParseInt(string(digits[1]), 10, 64)
	if err != nil {
		panic(err)
	}
	items[2], err = strconv.ParseInt(string(digits[2]), 10, 64)
	if err != nil {
		panic(err)
	}
	items[3], err = strconv.ParseInt(digits[3:], 10, 64)
	if err != nil {
		panic(err)
	}
	return items
}

func setValue(m mode, value int64, position int64, values []int64) {
	switch m {
	case immediateMode:
		values[position] = value
	case positionMode:
		values[values[position]] = value
	default:
		panic(fmt.Sprintf("mode %d not found", m))
	}
}

func getValue(m mode, position int64, values []int64) int64 {
	switch m {
	case immediateMode:
		return values[position]
	case positionMode:
		return values[values[position]]
	default:
		panic(fmt.Sprintf("mode %d not found", m))
	}
}

func newState(opCodePos int64, values []int64) []int64 {
	opcode := getOpcode(values[opCodePos])
	switch operator(opcode[3]) {
	case opAddition:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		setValue(mode(opcode[0]), pos1Val+pos2Val, opCodePos+3, values)
		return newState(opCodePos+4, values)
	case opMultiplier:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		setValue(mode(opcode[0]), pos1Val*pos2Val, opCodePos+3, values)
		return newState(opCodePos+4, values)
	case opInput:
		values[values[opCodePos+1]] = inputVal
		return newState(opCodePos+2, values)
	case opOutput:
		log.Printf("%d", values[values[opCodePos+1]])
		return newState(opCodePos+2, values)
	case opJumpTrue:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		if pos1Val != 0 {
			return newState(getValue(mode(opcode[1]), opCodePos+2, values), values)
		}
		return newState(opCodePos+3, values)
	case opJumpFalse:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		if pos1Val == 0 {
			return newState(getValue(mode(opcode[1]), opCodePos+2, values), values)
		}
		return newState(opCodePos+3, values)
	case lessThan:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		val := int64(0)
		if pos1Val < pos2Val {
			val = 1
		}
		setValue(mode(opcode[0]), val, opCodePos+3, values)
		return newState(opCodePos+4, values)
	case equals:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		val := int64(0)
		if pos1Val == pos2Val {
			val = 1
		}
		setValue(mode(opcode[0]), val, opCodePos+3, values)
		return newState(opCodePos+4, values)
	case opExit:
		return values
	default:
		panic(fmt.Sprintf("opcode not found: %d", opcode[3]))
	}
}

const inputVal = 5

func bruteForce(values []int64, pos1 int64, pos2 int64) {
	v := deepcopy.Copy(values).([]int64)
	newState(0, v)
}

func main() {
	values := readLines()
	bruteForce(values, 12, 0)
}
