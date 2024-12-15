package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	log.SetFlags(log.Lshortfile)

	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	var grid [][]byte
	var grid2 [][]byte
	var instrs []byte
	isGrid := true
	foundRobot := false
	var robot Vec2
	var robot2 Vec2

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			isGrid = false
			continue
		}
		if isGrid {
			line := []byte(scanner.Text())
			if !foundRobot {
				x := bytes.IndexByte(line, '@')
				if x >= 0 {
					foundRobot = true
					robot = Vec2{x: x, y: len(grid)}
					robot2 = Vec2{x: 2 * x, y: len(grid)}
					line[x] = '.'
				}
			}
			grid = append(grid, line)
			line2 := make([]byte, 0, len(line)*2)
			for _, i := range line {
				switch i {
				case '#':
					line2 = append(line2, '#', '#')
				case 'O':
					line2 = append(line2, '[', ']')
				case '.':
					line2 = append(line2, '.', '.')
				default:
					log.Fatalln("Unexpected grid char")
				}
			}
			grid2 = append(grid2, line2)
			continue
		}
		instrs = append(instrs, scanner.Bytes()...)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	for _, i := range instrs {
		sim(grid, &robot, i)
		sim2(grid2, &robot2, i)
	}

	sum := 0
	for y, i := range grid {
		for x, b := range i {
			if b == 'O' {
				sum += 100*y + x
			}
		}
	}
	fmt.Println("Part 1:", sum)

	sum = 0
	for y, i := range grid2 {
		for x, b := range i {
			if b == '[' {
				sum += 100*y + x
			}
		}
	}
	fmt.Println("Part 2:", sum)
}

type (
	Vec2 struct {
		x, y int
	}
)

func sim2(grid [][]byte, robot *Vec2, instr byte) {
	delta := deltaFromInstr(instr)
	next := addVec2(*robot, delta)
	if !canMoveBox(grid, next, delta) {
		return
	}
	moveBox(grid, next, delta)
	*robot = next
}

func moveBox(grid [][]byte, pos Vec2, delta Vec2) {
	isVertical := delta.y != 0
	next := addVec2(pos, delta)
	switch grid[pos.y][pos.x] {
	case '.':
		return
	case '#':
		return
	case '[':
		moveBox(grid, next, delta)
		if isVertical {
			side := next
			side.x++
			moveBox(grid, side, delta)
			grid[pos.y][side.x] = '.'
			grid[side.y][side.x] = ']'
		}
		grid[pos.y][pos.x] = '.'
		grid[next.y][next.x] = '['
	case ']':
		moveBox(grid, next, delta)
		if isVertical {
			side := next
			side.x--
			moveBox(grid, side, delta)
			grid[pos.y][side.x] = '.'
			grid[side.y][side.x] = '['
		}
		grid[pos.y][pos.x] = '.'
		grid[next.y][next.x] = ']'
	default:
		log.Fatalln("Unexpected grid char")
		return
	}
}

func canMoveBox(grid [][]byte, pos Vec2, delta Vec2) bool {
	isHorizontal := delta.x != 0
	next := addVec2(pos, delta)
	switch grid[pos.y][pos.x] {
	case '.':
		return true
	case '#':
		return false
	case '[':
		side := next
		side.x++
		return canMoveBox(grid, next, delta) && (isHorizontal || canMoveBox(grid, side, delta))
	case ']':
		side := next
		side.x--
		return canMoveBox(grid, next, delta) && (isHorizontal || canMoveBox(grid, side, delta))
	default:
		log.Fatalln("Unexpected grid char")
		return false
	}
}

func sim(grid [][]byte, robot *Vec2, instr byte) {
	delta := deltaFromInstr(instr)
	next := addVec2(*robot, delta)
	after := next
loop:
	for {
		switch grid[after.y][after.x] {
		case '.':
			break loop
		case '#':
			return
		case 'O':
			after = addVec2(after, delta)
		default:
			log.Fatalln("Unexpected grid char")
			return
		}
	}
	*robot = next
	grid[next.y][next.x] = '.'
	if after != next {
		grid[after.y][after.x] = 'O'
	}
}

func addVec2(a, b Vec2) Vec2 {
	return Vec2{x: a.x + b.x, y: a.y + b.y}
}

func deltaFromInstr(b byte) Vec2 {
	switch b {
	case '^':
		return Vec2{x: 0, y: -1}
	case '>':
		return Vec2{x: 1, y: 0}
	case 'v':
		return Vec2{x: 0, y: 1}
	case '<':
		return Vec2{x: -1, y: 0}
	default:
		log.Fatalln("Invalid instr")
		return Vec2{}
	}
}
