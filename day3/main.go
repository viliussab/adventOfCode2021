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

type ColumnOccurances struct {
	zero int
	one  int
}

func readFile(filename string) [][]int {
	var filePtr, err = os.Open("./input.txt")
	check(err)

	var scanner = bufio.NewScanner(filePtr)

	scanner.Split(bufio.ScanLines)

	var diagnostics [][]int

	for scanner.Scan() {
		var cmdLine []int
		lineStr := scanner.Text()
		chars := []rune(lineStr)
		for _, char := range chars {
			charStr := string(char)
			charInt, err := strconv.ParseInt(charStr, 10, 32)
			check(err)
			cmdLine = append(cmdLine, int(charInt))
		}

		diagnostics = append(diagnostics, cmdLine)
	}

	return diagnostics
}

func powerOfTwo(exponent int) int {
	return 1 << (exponent - 1)
}

func main() {
	var diagnostics = readFile("./input.txt")

	var colLen = len(diagnostics)
	var rowLen = len(diagnostics[0])

	var mostCommonBits []int

	for i := 0; i < rowLen; i++ {
		var colOccurances = ColumnOccurances{0, 0}

		for j := 0; j < colLen; j++ {

			if diagnostics[j][i] == 0 {
				colOccurances.zero++
			} else {
				colOccurances.one++
			}
		}

		var mostCommonBit int

		if colOccurances.zero > colOccurances.one {
			mostCommonBit = 0
		} else {
			mostCommonBit = 1
		}

		mostCommonBits = append(mostCommonBits, mostCommonBit)
	}

	var gamma int
	var epsilon int

	for i, bit := range mostCommonBits {
		exponent := (rowLen) - i

		if bit == 1 {
			gamma += powerOfTwo(exponent)
		} else {
			epsilon += powerOfTwo(exponent)
		}

	}

	fmt.Println(gamma * epsilon)

	getLSMostInteresting := func(colOccurances ColumnOccurances) int {
		if colOccurances.zero > colOccurances.one {
			return 0
		} else {
			return 1
		}
	}

	var oxygenRating = calcLifeSupportRating(diagnostics, getLSMostInteresting)

	getCO2MostInteresting := func(colOccurances ColumnOccurances) int {
		if colOccurances.zero > colOccurances.one {
			return 1
		} else {
			return 0
		}
	}

	var co2Rating = calcLifeSupportRating(diagnostics, getCO2MostInteresting)

	fmt.Println("Oxygen Rating", oxygenRating)
	fmt.Println("CO2 Rating", co2Rating)

	fmt.Println(oxygenRating * co2Rating)
}

type getMostInterestingBit func(colOccurances ColumnOccurances) int

func calcLifeSupportRating(diagnostics [][]int, getMostInterestingBit getMostInterestingBit) int {

	var remainingRows [][]int = diagnostics
	var acceptedRow []int
	var rowLen = len(diagnostics[0])

	for i := 0; i < rowLen; i++ {

		var columnLen = len(remainingRows)

		var colOccurances = ColumnOccurances{0, 0}

		for j := 0; j < columnLen; j++ {

			if remainingRows[j][i] == 0 {
				colOccurances.zero++
			} else {
				colOccurances.one++
			}
		}

		var bitOfInterest = getMostInterestingBit(colOccurances)

		var acceptedRows [][]int

		for j := 0; j < columnLen; j++ {
			if remainingRows[j][i] == bitOfInterest {
				acceptedRows = append(acceptedRows, remainingRows[j])
			}
		}

		remainingRows = acceptedRows

		if len(remainingRows) == 1 {
			acceptedRow = remainingRows[0]
			break
		}
	}

	var res int
	for i, bit := range acceptedRow {
		exponent := (rowLen) - i

		if bit == 1 {
			res += powerOfTwo(exponent)
		}
	}

	return res
}
