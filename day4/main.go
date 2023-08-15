package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type BoardNumber struct {
	value   int
	occured bool
}

type PlayerBoard [5][5]BoardNumber

type Bingo struct {
	rolls  []int
	boards []PlayerBoard
}

func parseTurns(numLine string) []int {
	var buf []int
	var nums = strings.Split(numLine, ",")

	for _, num := range nums {
		intNum, err := strconv.ParseInt(num, 10, 32)
		check(err)
		buf = append(buf, int(intNum))
	}

	return buf
}

func parseBoard(boardLines [5]string) PlayerBoard {
	var grid PlayerBoard

	var parseLineNums = func(lineNums [5]string) [5]BoardNumber {
		var numbers [5]BoardNumber
		for i, strNum := range lineNums {
			intNum, err := strconv.ParseInt(strNum, 10, 32)
			check(err)
			numbers[i] = BoardNumber{
				value:   int(intNum),
				occured: false}
		}

		return numbers
	}

	for i, line := range boardLines {
		var lineNums = strings.Fields(line)
		var strNums = [5]string{}
		copy(strNums[:], lineNums)
		var boardNums = parseLineNums(strNums)

		grid[i] = boardNums
	}

	return grid
}

func readFile(filename string) Bingo {
	var filePtr, err = os.Open("./input.txt")
	check(err)

	var scanner = bufio.NewScanner(filePtr)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	var turnsStr = scanner.Text()
	var turns = parseTurns(turnsStr)

	var scanBoard = func(scanner *bufio.Scanner) [5]string {
		var board [5]string

		for i := 0; i < 5; i++ {
			scanner.Scan()
			board[i] = scanner.Text()
		}

		return board
	}

	var boards []PlayerBoard

	for scanner.Scan() {
		var rawBoard = scanBoard(scanner)
		var parsedBoard = parseBoard(rawBoard)
		boards = append(boards, parsedBoard)
	}

	return Bingo{
		rolls:  turns,
		boards: boards,
	}
}

func assignRolls(boards *[]PlayerBoard, roll int) {
	for k, board := range *boards {
		result := findNumberIndex(board, roll)
		if !result.isFound {
			continue
		}

		i, j := result.i, result.j

		board[i][j].occured = true

		(*boards)[k] = board
	}
}

type FindResult struct {
	isFound bool
	i       int
	j       int
}

func findNumberIndex(board PlayerBoard, number int) FindResult {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if board[i][j].value == number {
				return FindResult{
					isFound: true,
					i:       i,
					j:       j,
				}
			}
		}
	}

	return FindResult{
		isFound: false,
		i:       -1,
		j:       -1,
	}
}

type IsWonResult struct {
	isWin      bool
	winIndexes []int
}

func isBingo(boards *[]PlayerBoard) IsWonResult {
	var indexes []int

	for i, board := range *boards {
		if hasBingo(board) {
			indexes = append(indexes, i)
		}
	}

	return IsWonResult{
		isWin:      len(indexes) > 0,
		winIndexes: indexes,
	}
}

func hasBingo(board PlayerBoard) bool {
	for i := 0; i < 5; i++ {
		var row = true
		var col = true
		for j := 0; j < 5; j++ {
			if !board[i][j].occured {
				row = false
			}

			if !board[j][i].occured {
				col = false
			}
		}

		if row || col {
			return true
		}
	}

	return false
}

func findUnoccupiedSum(board PlayerBoard) int {
	var sum = 0

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if !board[i][j].occured {
				sum += board[i][j].value
			}
		}
	}

	return sum
}

func main() {
	variantA()
	variantB()
}

func variantB() {
	var bingo = readFile("./input.txt")

	var remainingIndexes []int

	for i := range bingo.boards {
		remainingIndexes = append(remainingIndexes, i)
	}

	var lastWinners []PlayerBoard
	var completedRoll int = -1

	for _, roll := range bingo.rolls {
		assignRolls(&bingo.boards, roll)

		var result = isBingo(&bingo.boards)
		if len(result.winIndexes) == len(bingo.boards) {
			for _, loseIndex := range remainingIndexes {
				lastWinners = append(lastWinners, bingo.boards[loseIndex])
			}
			completedRoll = roll
			break
		}

		var left []int

		for _, remainingIndex := range remainingIndexes {
			var found = false
			for _, windex := range result.winIndexes {
				if windex == remainingIndex {
					found = true
				}
			}

			if !found {
				left = append(left, remainingIndex)
			}
		}

		remainingIndexes = left
	}

	if len(lastWinners) == 0 || completedRoll == -1 {
		panic("Not all players finished")
	}

	var oneOfLast = lastWinners[0]
	var sum = findUnoccupiedSum(oneOfLast)

	println(sum * completedRoll)
}

func variantA() {
	var bingo = readFile("./input.txt")

	var winningPlayers []PlayerBoard
	var winningRoll int = -1

	for _, roll := range bingo.rolls {

		assignRolls(&bingo.boards, roll)

		var result = isBingo(&bingo.boards)
		if result.isWin {
			winningRoll = roll
			for _, winIndex := range result.winIndexes {
				winningPlayers = append(winningPlayers, bingo.boards[winIndex])
			}
			break
		}
	}

	if len(winningPlayers) == 0 || winningRoll == -1 {
		panic("No winning players found")
	}

	var firstWinner = winningPlayers[0]
	var sum = findUnoccupiedSum(firstWinner)

	println(sum * winningRoll)
}
