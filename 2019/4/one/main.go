package main

import (
	"log"
	"strconv"
)

var start = 359282
var end = 820401


func isValid(number string) bool {
	return isIncreasing(number)
}

func isIncreasing(number string) bool {
	previous := uint8(0)
	hasDouble := uint8(8)
	for i := 0; i < len(number); i++ {
		if previous > number[i] {
			return false
		}
		previous = number[i]
	}

	return hasDouble > 0
}

func main() {
	correct := 0
	for i := start; i < end + 1; i++ {
		number := strconv.FormatInt(int64(i), 10)
		if isValid(number) {
			correct++
		}
	}
	log.Printf("%d", correct)
}
