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

func getSources() map[string]Target {
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
		target, ok := targets[o[1]]
		if !ok {
			target = make([]string, 0)
		}
		target = append(target, o[0])
		targets[o[1]] = target
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

func getPath(source string, target string, targets map[string]Target) []string {
	path := make([]string, 0)
	if source == target {
		return path
	}
	a, ok := targets[source]
	if len(a) == 0 && !ok {
		panic("Not found")
	}
	path = append(path, a[0])
	return append(path, getPath(a[0], target, targets)...)
}

const centerOfMass = "COM"

func main() {
	//values := getTargets()
	//log.Printf("orbits: %d", determineOrbits(centerOfMass, 0, values))

    values := getSources()
    youPath := getPath("YOU", centerOfMass, values)
    targetPath := getPath("SAN", centerOfMass, values)

    log.Printf("%+v %+v", youPath, targetPath)
    same := ""
    youIdx := 0
    targetIdx := 0
    smallest := 10000000
    for idx, you := range youPath {
    	for idx_2, san := range targetPath {
    		if you == san {
    			if idx + idx_2 < smallest {
					youIdx = idx
					targetIdx = idx_2
					same = san
					smallest = idx + idx_2
				}
			}
		}
	}
	log.Printf("same %s %d %d", same, youIdx, targetIdx)

}

