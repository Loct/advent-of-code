package main

import (
	"fmt"
	"github.com/mohae/deepcopy"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
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
		values[values[opCodePos+1]] = <-inputs
		return newState(opCodePos+2, values, inputs, output)
	case opOutput:
		//log.Printf("%d", values[values[opCodePos+1]])
		output <- values[values[opCodePos+1]]
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
	case opExit:
		halted++
		return values
	default:
		panic(fmt.Sprintf("opcode not found: %d", opcode[3]))
	}
}

var halted int64

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

type amp struct {
	input  chan int64
	output chan int64
	values []int64
}

func phase(values []int64, phases []int64) {
	halted = 0
	ampsLock := sync.RWMutex{}
	amps := make([]amp, len(phases))
	amps[0].input = make(chan int64, 4)
	amps[0].input <- phases[0]
	amps[0].input <- 0
	amps[0].output = make(chan int64, 4)

	amps[1].input = amps[0].output
	amps[1].input <- phases[1]

	amps[1].output = make(chan int64, 4)

	amps[2].input = amps[1].output
	amps[2].input <- phases[2]

	amps[2].output = make(chan int64, 4)

	amps[3].input = amps[2].output
	amps[3].input <- phases[3]

	amps[3].output = make(chan int64, 4)

	amps[4].input = amps[3].output
	amps[4].input <- phases[4]
	amps[4].output = amps[0].input

	for idx, _ := range phases {
		ampsLock.Lock()
		amp := amps[idx]
		ampsLock.Unlock()
		v := deepcopy.Copy(values).([]int64)
		amp.values = v
		go func() {
			newState(0, amp.values, amp.input, amp.output)
		}()
	}

	for halted != 5 {
		//log.Printf("busy")
	}

	for _, amp := range amps {
		select {
		case m := <-amp.output:
			if m > maxValue {
				log.Printf("last %d %+v", m, phases)
				maxValue = m
			}
		default:
			break
		}
	}
}

var maxValue = int64(-1)

func main() {
	ph = make([][]int64, 0)
	values := readLines()
	//ph = append(ph, []int64{9,8,7,6,5})
	heapPermutation([]int64{9, 8, 7, 6, 5}, 5, 0)
	for _, p := range ph {
		log.Printf("%d", maxValue)
		phase(values, p)
	}
}
