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
const opLessThan operator = 7
const opEquals operator = 8
const opRelativeBaseUpdate operator = 9

const opExit operator = 99

type mode int64

const positionMode mode = 0
const immediateMode mode = 1
const relativeMode mode = 2

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
	case relativeMode:
		values[relativeBase+values[position]] = value
	default:
		panic(fmt.Sprintf("mode %d not found", m))
	}
}

var relativeBase int64

func getValue(m mode, position int64, values []int64) int64 {
	switch m {
	case immediateMode:
		return values[position]
	case positionMode:
		return values[values[position]]
	case relativeMode:
		return values[relativeBase+values[position]]
	default:
		panic(fmt.Sprintf("mode %d not found", m))
	}
}

func newState(opCodePos int64, values []int64, inputs chan int64, output chan int64) []int64 {
	opcode := getOpcode(values[opCodePos])
	switch operator(opcode[3]) {
	case opAddition:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		setValue(mode(opcode[0]), pos1Val+pos2Val, opCodePos+3, values)
		return newState(opCodePos+4, values, inputs, output)
	case opMultiplier:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		setValue(mode(opcode[0]), pos1Val*pos2Val, opCodePos+3, values)
		return newState(opCodePos+4, values, inputs, output)
	case opInput:
		v := <-inputs
		setValue(mode(opcode[2]), v, opCodePos+1, values)
		return newState(opCodePos+2, values, inputs, output)
	case opOutput:
		v := getValue(mode(opcode[2]), opCodePos+1, values)
		output <- v
		return newState(opCodePos+2, values, inputs, output)
	case opJumpTrue:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		if pos1Val != 0 {
			return newState(getValue(mode(opcode[1]), opCodePos+2, values), values, inputs, output)
		}
		return newState(opCodePos+3, values, inputs, output)
	case opJumpFalse:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		if pos1Val == 0 {
			return newState(getValue(mode(opcode[1]), opCodePos+2, values), values, inputs, output)
		}
		return newState(opCodePos+3, values, inputs, output)
	case opLessThan:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		val := int64(0)
		if pos1Val < pos2Val {
			val = 1
		}
		setValue(mode(opcode[0]), val, opCodePos+3, values)
		return newState(opCodePos+4, values, inputs, output)
	case opEquals:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		val := int64(0)
		if pos1Val == pos2Val {
			val = 1
		}
		setValue(mode(opcode[0]), val, opCodePos+3, values)
		return newState(opCodePos+4, values, inputs, output)
	case opRelativeBaseUpdate:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		relativeBase = relativeBase + pos1Val
		return newState(opCodePos+2, values, inputs, output)
	case opExit:
		halted++
		return values
	default:
		panic(fmt.Sprintf("opcode not found: %d", opcode[3]))
	}
}

var halted int64

type amp struct {
	input  chan int64
	output chan int64
	values []int64
}

func run(values []int64) {
	halted = 0
	amps := make([]amp, 1)
	amps[0].input = make(chan int64, 2)
	amps[0].input <- 1
	amps[0].output = make(chan int64, 100)

	for len(values) < 5000 {
		values = append(values, 0)
	}
	v := deepcopy.Copy(values).([]int64)
	newState(0, v, amps[0].input, amps[0].output)

	done := false
	for {
		select {
		case m := <-amps[0].output:
			log.Printf("%d", m)
		default:
			done = true
			break
		}
		if done {
			break
		}
	}
}

func main() {
	values := readLines()
	run(values)
}
