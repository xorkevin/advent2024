package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/xorkevin/advent2024/bitset"
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
	needStart := true
	var start Pos

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := []byte(scanner.Text())
		if needStart {
			if col := bytes.IndexByte(row, '^'); col >= 0 {
				needStart = false
				start = Pos{x: col, y: len(grid)}
			}
		}
		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	visited := sim(grid, start)
	fmt.Println("Part 1:", len(visited)+1)

	count2 := 0
	h := len(grid)
	w := len(grid[0])
	seen := bitset.New(w * h * 4)
	for _, i := range visited {
		if simLoop(grid, start, i, seen) {
			count2++
		}
	}
	fmt.Println("Part 2:", count2)
}

func simLoop(grid [][]byte, pos Pos, obstruct Pos, seen *bitset.BitSet) bool {
	h := len(grid)
	w := len(grid[0])

	delta := Pos{x: 0, y: -1}
	seen.Reset()
	for {
		next := addPos(pos, delta)
		if outBounds(next, w, h) {
			return false
		}
		if next == obstruct || grid[next.y][next.x] == '#' {
			delta = turn(delta)
			if !seen.Insert(getStateID(pos, w, delta)) {
				return true
			}
			continue
		}
		pos = next
	}
}

func sim(grid [][]byte, pos Pos) []Pos {
	h := len(grid)
	w := len(grid[0])

	var p []Pos
	delta := Pos{x: 0, y: -1}
	seen := bitset.New(w * h)
	seen.Insert(getPosID(pos, w))
	for {
		next := addPos(pos, delta)
		if outBounds(next, w, h) {
			break
		}
		if grid[next.y][next.x] == '#' {
			delta = turn(delta)
			continue
		}
		pos = next
		if seen.Insert(getPosID(pos, w)) {
			p = append(p, pos)
		}
	}
	return p
}

type (
	Pos struct {
		x, y int
	}
)

func getPosID(pos Pos, w int) int {
	return pos.y*w + pos.x
}

func getDeltaID(delta Pos) int {
	if delta.y < 0 {
		return 0
	}
	if delta.x > 0 {
		return 1
	}
	if delta.y > 0 {
		return 2
	}
	return 3
}

func getStateID(pos Pos, w int, delta Pos) int {
	return getPosID(pos, w)*4 + getDeltaID(delta)
}

func addPos(a, b Pos) Pos {
	return Pos{x: a.x + b.x, y: a.y + b.y}
}

func outBounds(pos Pos, w, h int) bool {
	return pos.x < 0 || pos.y < 0 || pos.x >= w || pos.y >= h
}

func turn(delta Pos) Pos {
	if delta.y < 0 {
		return Pos{x: 1, y: 0}
	}
	if delta.x > 0 {
		return Pos{x: 0, y: 1}
	}
	if delta.y > 0 {
		return Pos{x: -1, y: 0}
	}
	return Pos{x: 0, y: -1}
}
