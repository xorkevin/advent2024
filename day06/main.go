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

	count1, visited := sim(grid, start)
	fmt.Println("Part 1:", count1)

	count2 := 0
	h := len(grid)
	w := len(grid[0])
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if grid[i][j] != '.' {
				continue
			}
			pos := Pos{x: j, y: i}
			if !visited[getPosID(pos, w)] {
				continue
			}
			if simLoop(grid, start, pos) {
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
	seen := make([]byte, w*h)
	seen[getPosID(pos, w)] |= getDeltaID(delta)
	for {
		next := addPos(pos, delta)
		if outBounds(next, w, h) {
			break
		}
		if next == obstruct || grid[next.y][next.x] == '#' {
			delta = turn(delta)
			id := getPosID(pos, w)
			deltaID := getDeltaID(delta)
			if seen[id]&deltaID != 0 {
				return true
			}
			seen[id] |= deltaID
			continue
		}
		pos = next
		id := getPosID(pos, w)
		deltaID := getDeltaID(delta)
		if seen[id]&deltaID != 0 {
			return true
		}
		seen[id] |= deltaID
	}
	return false
}

func sim(grid [][]byte, pos Pos) (int, []bool) {
	h := len(grid)
	w := len(grid[0])

	delta := Pos{x: 0, y: -1}
	seen := make([]bool, w*h)
	seen[getPosID(pos, w)] = true
	count := 1
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
		id := getPosID(pos, w)
		if !seen[id] {
			count++
		}
		seen[id] = true
	}
	return count, seen
}

type (
	Pos struct {
		x, y int
	}
)

func getPosID(pos Pos, w int) int {
	return pos.y*w + pos.x
}

func getDeltaID(delta Pos) byte {
	if delta.y < 0 {
		return 1
	}
	if delta.x > 0 {
		return 1 << 1
	}
	if delta.y > 0 {
		return 1 << 2
	}
	return 1 << 3
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
