package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(filename string) []int {
	var filePtr, err = os.Open("./input.txt")
	check(err)

	var scanner = bufio.NewScanner(filePtr)

	scanner.Split(bufio.ScanLines)

	var lines []int

	for scanner.Scan() {
		lineStr := scanner.Text()
		lineInt, err := strconv.ParseInt(lineStr, 10, 32)
		check(err)
		lines = append(lines, int(lineInt))
	}

	return lines
}

func calcLatestThree(measurements []int, i int) int {
	if i < 2 {
		panic(fmt.Sprintf("i=%v is out of bounds", i))
	}

	return measurements[i-2] + measurements[i-1] + measurements[i]
}

func main() {
	var measurements = readFile("./input.txt")

	increases := 0
	decreases := 0

	for i, _ := range measurements {
		if i < 3 {
			continue
		}

		curr := calcLatestThree(measurements, i)
		prev := calcLatestThree(measurements, i-1)

		if curr > prev {
			increases += 1
		} else if curr < prev {
			decreases += 1
		}
	}

	fmt.Println(increases)
}
