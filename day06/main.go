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
	needStart := true
	var start Pos

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := []byte(scanner.Text())
		if needStart {
			if col := bytes.IndexByte(row, '^'); col >= 0 {
				start = Pos{x: col, y: len(grid)}
			}
		}
		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	count1 := sim(grid, start)
	fmt.Println("Part 1:", count1)

	count2 := 0
	h := len(grid)
	w := len(grid[0])
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if grid[i][j] != '.' {
				continue
			}
			if simLoop(grid, start, Pos{x: j, y: i}) {
				count2++
			}
		}
	}
	fmt.Println("Part 2:", count2)
}

func simLoop(grid [][]byte, pos Pos, obstruct Pos) bool {
	h := len(grid)
	w := len(grid[0])

	delta := Pos{x: 0, y: -1}
	seen := map[State]struct{}{}
	seen[State{pos: pos, delta: delta}] = struct{}{}
	for {
		next := addPos(pos, delta)
		if outBounds(next, w, h) {
			break
		}
		if next == obstruct || grid[next.y][next.x] == '#' {
			delta = turn(delta)
			s := State{pos: pos, delta: delta}
			if _, ok := seen[s]; ok {
				return true
			}
			seen[s] = struct{}{}
			continue
		}
		pos = next
		s := State{pos: pos, delta: delta}
		if _, ok := seen[s]; ok {
			return true
		}
		seen[s] = struct{}{}
	}
	return false
}

type (
	State struct {
		pos, delta Pos
	}
)

func sim(grid [][]byte, pos Pos) int {
	h := len(grid)
	w := len(grid[0])

	delta := Pos{x: 0, y: -1}
	seen := map[Pos]struct{}{}
	seen[pos] = struct{}{}
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
		seen[pos] = struct{}{}
	}
	return len(seen)
}

type (
	Pos struct {
		x, y int
	}
)

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
