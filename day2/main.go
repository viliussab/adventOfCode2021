package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Direction int

const (
	Forward = iota
	Up
	Down
)

var directions = map[string]Direction{
	"forward": Forward,
	"up":      Up,
	"down":    Down,
}

var directionsAsString = map[Direction]string{
	Forward: "forward",
	Up:      "up",
	Down:    "down",
}

type SubmarineCommand struct {
	direction Direction
	distance  int
}

type Submarine struct {
	forward int
	depth   int
	aim     int
}

func readFile(filename string) []SubmarineCommand {
	var filePtr, err = os.Open("./input.txt")
	check(err)

	var scanner = bufio.NewScanner(filePtr)

	scanner.Split(bufio.ScanLines)

	var cmds []SubmarineCommand

	for scanner.Scan() {
		lineStr := scanner.Text()
		lineParts := strings.Split(lineStr, " ")
		dirStr, distStr := lineParts[0], lineParts[1]

		dir := directions[dirStr]
		dist, err := strconv.ParseInt(distStr, 10, 32)
		check(err)
		cmd := SubmarineCommand{
			direction: dir,
			distance:  int(dist),
		}

		cmds = append(cmds, cmd)
	}

	return cmds
}

func (pos *Submarine) execudeCommand(cmd SubmarineCommand) {
	switch cmd.direction {
	case Forward:
		pos.forward += cmd.distance
		pos.depth += pos.aim * cmd.distance
	case Up:
		pos.aim -= cmd.distance
	case Down:
		pos.aim += cmd.distance
	}
}

func main() {
	var commands = readFile("./input.txt")

	position := Submarine{
		forward: 0,
		depth:   0,
		aim:     0,
	}

	for _, cmd := range commands {
		position.execudeCommand(cmd)
	}

	fmt.Println(position.forward * position.depth)
}
