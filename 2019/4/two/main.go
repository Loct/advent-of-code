package main

import (
	"log"
	"strconv"
)

var start = 359282
var end = 820401


func isValid(number string) bool {
	log.Printf("%v %v", isIncreasing(number), hasAtleastOneDouble(number))
	return isIncreasing(number) && hasAtleastOneDouble(number)
}

func hasAtleastOneDouble(number string) bool {
	items := make(map[uint8]int)
	previous := uint8(0)
	for i := 0; i < len(number); i++ {
		if previous == 0 {
			previous = number[i]
		}

		if previous != number[i] {
			previous = number[i]
		}
		items[previous]++

	}
	for _, value := range items {
		if value == 2 {
			return true
		}
	}
	return false
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
