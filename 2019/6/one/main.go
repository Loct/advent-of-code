package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Target []string

func getTargets() map[string]Target {
	f, err := os.Open("../data")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic("err")
	}
	values := strings.Split(string(data), "\n")
	targets := make(map[string]Target, 0)
	for _, value := range values {
		o := strings.Split(value, ")")
		if len(o) != 2 {
			continue
		}
		target, ok := targets[o[0]]
		if !ok {
			target = make([]string, 0)
		}
		target = append(target, o[1])
		targets[o[0]] = target
	}
	return targets
}

func determineOrbits(origin string, currentOrbits int, targets map[string]Target) int {
	target, _ := targets[origin]
	if len(target) == 0 {
		return currentOrbits
	}
	additionalOrbits := 0
	for _, t := range target {
		additionalOrbits += determineOrbits(t, currentOrbits + 1, targets)
	}
	return  additionalOrbits + currentOrbits
}

const centerOfMass = "COM"

func main() {
	values := getTargets()
	log.Printf("orbits: %d", determineOrbits(centerOfMass, 0, values))
}
