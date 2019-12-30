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

var currentOutput int64
func newState(opCodePos int64, values []int64, inputs []int64) []int64 {
	opcode := getOpcode(values[opCodePos])
	switch operator(opcode[3]) {
	case opAddition:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		setValue(mode(opcode[0]), pos1Val+pos2Val, opCodePos+3, values)
		return newState(opCodePos+4, values, inputs)
	case opMultiplier:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		setValue(mode(opcode[0]), pos1Val*pos2Val, opCodePos+3, values)
		return newState(opCodePos+4, values, inputs)
	case opInput:
		values[values[opCodePos+1]] = inputs[0]
		inputs = inputs[1:]
		return newState(opCodePos+2, values, inputs)
	case opOutput:
		currentOutput =  values[values[opCodePos+1]]
		return newState(opCodePos+2, values, inputs)
	case opJumpTrue:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		if pos1Val != 0 {
			return newState(getValue(mode(opcode[1]), opCodePos+2, values), values, inputs)
		}
		return newState(opCodePos+3, values, inputs)
	case opJumpFalse:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		if pos1Val == 0 {
			return newState(getValue(mode(opcode[1]), opCodePos+2, values), values, inputs)
		}
		return newState(opCodePos+3, values, inputs)
	case opLessThan:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		val := int64(0)
		if pos1Val < pos2Val {
			val = 1
		}
		setValue(mode(opcode[0]), val, opCodePos+3, values)
		return newState(opCodePos+4, values, inputs)
	case opEquals:
		pos1Val := getValue(mode(opcode[2]), opCodePos+1, values)
		pos2Val := getValue(mode(opcode[1]), opCodePos+2, values)
		val := int64(0)
		if pos1Val == pos2Val {
			val = 1
		}
		setValue(mode(opcode[0]), val, opCodePos+3, values)
		return newState(opCodePos+4, values, inputs)
	case opExit:
		return values
	default:
		panic(fmt.Sprintf("opcode not found: %d", opcode[3]))
	}
}

var ph [][]int64
func heapPermutation(a []int64, size int64, n int64) []int64 {
	// if size becomes 1 then prints the obtained
	// permutation
	if size == 1 {
		//log.Printf("%+v", a)
		ph = append(ph, deepcopy.Copy(a).([]int64))
		return a
	}

	for i := int64(0); i < size; i++ {
		a = heapPermutation(a, size-1, n)
		// if size is odd, swap first and last
		// element
		if size%2 == 1 {
			temp := a[0]
			a[0] = a[size-1]
			a[size-1] = temp
		} else {
			temp := a[i]
			a[i] = a[size-1]
			a[size-1] = temp
		}
	}
    return a
}

var maxValue = int64(-1)
func phase(values []int64, phases []int64) {
	for _, phase := range phases {
		v := deepcopy.Copy(values).([]int64)
		newState(0, v, []int64{phase, currentOutput})
	}
	if currentOutput > maxValue {
		log.Printf("%d, %v", currentOutput, phases)
		maxValue = currentOutput
	}
	currentOutput = 0
}

func main() {
	ph = make([][]int64, 0)
	values := readLines()
	heapPermutation([]int64{0,1,2,3,4}, 5, 0)
	for _, p := range ph {
		phase(values, p)
	}
}
