package main

import (
	"bufio"
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	h := len(grid)
	w := len(grid[0])

	s1 := 0
	s2 := 0
	seen := bitset.New(w * h)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if grid[i][j] != '0' {
				continue
			}
			seen.Reset()
			r1, r2 := search(grid, Pos{x: j, y: i}, 0, seen)
			s1 += r1
			s2 += r2
		}
	}
	fmt.Println("Part 1:", s1)
	fmt.Println("Part 2:", s2)
}

type (
	Pos struct {
		x, y int
	}
)

func search(grid [][]byte, pos Pos, elevation int, seen *bitset.BitSet) (int, int) {
	if elevation == 9 {
		return 1, 1
	}
	h := len(grid)
	w := len(grid[0])
	bounds := Pos{
		x: w,
		y: h,
	}
	s1 := 0
	s2 := 0
	for _, a := range [4]Pos{
		{x: pos.x, y: pos.y - 1},
		{x: pos.x + 1, y: pos.y},
		{x: pos.x, y: pos.y + 1},
		{x: pos.x - 1, y: pos.y},
	} {
		if outBounds(a, bounds) {
			continue
		}
		c := grid[a.y][a.x] - '0'
		if c > 9 || int(c) != elevation+1 {
			continue
		}
		id := getPosID(a, w)
		addedPoint := seen.Insert(id)
		r1, r2 := search(grid, a, elevation+1, seen)
		if addedPoint {
			s1 += r1
		}
		s2 += r2
	}
	return s1, s2
}

func getPosID(pos Pos, w int) int {
	return pos.y*w + pos.x
}

func outBounds(pos Pos, bounds Pos) bool {
	return pos.x < 0 || pos.y < 0 || pos.x >= bounds.x || pos.y >= bounds.y
}
