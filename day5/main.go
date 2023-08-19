package main

import (
	"bufio"
	"errors"
	"math"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type LineMovement struct {
	from Coordinate
	to   Coordinate
}

type Coordinate struct {
	x int
	y int
}

type CoordOccur struct {
	coord     Coordinate
	frequency int
}

func readFile(filename string) []LineMovement {
	var filePtr, err = os.Open(filename)
	check(err)

	var scanner = bufio.NewScanner(filePtr)
	scanner.Split(bufio.ScanLines)

	var parseCoordinate = func(cordStr string) Coordinate {
		var parts = strings.Split(cordStr, ",")
		x, err := strconv.ParseInt(parts[0], 10, 32)
		check(err)
		y, err := strconv.ParseInt(parts[1], 10, 32)
		check(err)

		return Coordinate{
			x: int(x),
			y: int(y),
		}
	}

	var movements []LineMovement

	for scanner.Scan() {
		var line = scanner.Text()

		var parts = strings.Split(line, " -> ")
		var from, to = parseCoordinate(parts[0]), parseCoordinate(parts[1])

		movements = append(movements, LineMovement{
			from: from,
			to:   to,
		})
	}

	return movements
}

type FilterResult struct {
	nonDiagonal []LineMovement
	diagonal    []LineMovement
}

func filterOutDiagonal(movements []LineMovement) FilterResult {
	var nonDiagonal []LineMovement
	var diagonal []LineMovement

	var calcDist1D = func(axisCoord int, axisCoordOther int) int {
		var floatDiff = math.Abs(float64(axisCoord - axisCoordOther))
		var diff = int(floatDiff)
		return diff
	}

	for _, mov := range movements {
		if mov.from.x == mov.to.x || mov.from.y == mov.to.y {
			nonDiagonal = append(nonDiagonal, mov)
		} else if calcDist1D(mov.from.x, mov.to.x) == calcDist1D(mov.from.y, mov.to.y) {
			diagonal = append(diagonal, mov)
		}
	}

	return FilterResult{
		nonDiagonal: nonDiagonal,
		diagonal:    diagonal,
	}
}

type TraverseDirection string

const (
	X TraverseDirection = "X"
	Y TraverseDirection = "Y"
)

func main() {
	var readResult = readFile("./input.txt")

	var filterRes = filterOutDiagonal(readResult)
	var diagonal, nonDiagonal = filterRes.diagonal, filterRes.nonDiagonal

	var coordinateMap = make(map[string]CoordOccur)
	var getDictionaryKey = func(x int, y int) string {
		return strconv.Itoa(x) + "," + strconv.Itoa(y)
	}
	var appendToMap = func(coord Coordinate) {
		var key = getDictionaryKey(coord.x, coord.y)

		x, found := coordinateMap[key]
		if found {
			x.frequency += 1
			coordinateMap[key] = x
		} else {
			coordinateMap[key] = CoordOccur{
				coord:     coord,
				frequency: 1,
			}
		}
	}

	var getMovementCoord = func(dir TraverseDirection, constCoord int, movCoord int) (Coordinate, error) {
		if dir == X {
			return Coordinate{
				x: movCoord,
				y: constCoord,
			}, nil
		}

		if dir == Y {
			return Coordinate{
				x: constCoord,
				y: movCoord,
			}, nil
		}

		return Coordinate{}, errors.New("not an enum value")
	}

	var getMinMax = func(val int, other int) (int, int) {
		if val < other {
			return val, other
		}

		return other, val
	}

	for _, mov := range nonDiagonal {

		var direction TraverseDirection
		var constantCoord int
		var movCoordFrom int
		var movCoordTo int

		if mov.from.x == mov.to.x {
			direction = Y
			constantCoord = mov.from.x
			movCoordFrom, movCoordTo = getMinMax(mov.from.y, mov.to.y)
		} else {
			direction = X
			constantCoord = mov.from.y
			movCoordFrom, movCoordTo = getMinMax(mov.from.x, mov.to.x)
		}

		for movCoord := movCoordFrom; movCoord <= movCoordTo; movCoord += 1 {
			var coord, err = getMovementCoord(direction, constantCoord, movCoord)
			check(err)
			appendToMap(coord)
		}
	}

	for _, mov := range diagonal {
		var xStep, yStep int

		if mov.from.x < mov.to.x {
			xStep = 1
		} else {
			xStep = -1
		}

		if mov.from.y < mov.to.y {
			yStep = 1
		} else {
			yStep = -1
		}
		var x, y = mov.from.x, mov.from.y
		for {
			appendToMap(Coordinate{x: x, y: y})

			x += xStep
			y += yStep

			var outOfBounds = func(i, n, step int) bool {
				var smaller, larger int
				if step > 0 {
					smaller, larger = i, n
				} else {
					smaller, larger = n, i
				}

				return smaller > larger
			}

			if outOfBounds(x, mov.to.x, xStep) || outOfBounds(y, mov.to.y, yStep) {
				break
			}
		}
	}

	var count = 0

	for _, v := range coordinateMap {
		if v.frequency >= 2 {
			count += 1
		}
	}

	println(count)
}
