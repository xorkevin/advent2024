package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

func main() {
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

	count1 := 0
	m1 := Matcher{
		target: "XMAS",
	}
	m2 := Matcher{
		target: "SAMX",
	}
	h := len(grid)
	w := len(grid[0])
	for i := 0; i < h; i++ {
		m1.Reset()
		m2.Reset()
		for j := 0; j < w; j++ {
			c := grid[i][j]
			if m1.Match(c) {
				count1++
			}
			if m2.Match(c) {
				count1++
			}
		}
	}
	for j := 0; j < w; j++ {
		m1.Reset()
		m2.Reset()
		for i := 0; i < h; i++ {
			c := grid[i][j]
			if m1.Match(c) {
				count1++
			}
			if m2.Match(c) {
				count1++
			}
		}
	}
	for j := 0; j < w; j++ {
		m1.Reset()
		m2.Reset()
		for k := 0; k < h && j+k < w; k++ {
			c := grid[k][j+k]
			if m1.Match(c) {
				count1++
			}
			if m2.Match(c) {
				count1++
			}
		}
	}
	for i := 1; i < h; i++ {
		m1.Reset()
		m2.Reset()
		for k := 0; i+k < h && k < w; k++ {
			c := grid[i+k][k]
			if m1.Match(c) {
				count1++
			}
			if m2.Match(c) {
				count1++
			}
		}
	}
	for j := 0; j < w; j++ {
		m1.Reset()
		m2.Reset()
		for k := 0; k < h && j-k >= 0; k++ {
			c := grid[k][j-k]
			if m1.Match(c) {
				count1++
			}
			if m2.Match(c) {
				count1++
			}
		}
	}
	for i := 1; i < h; i++ {
		m1.Reset()
		m2.Reset()
		for k := 0; i+k < h && w-k-1 >= 0; k++ {
			c := grid[i+k][w-k-1]
			if m1.Match(c) {
				count1++
			}
			if m2.Match(c) {
				count1++
			}
		}
	}
	fmt.Println("Part 1:", count1)

	count2 := 0
	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			e := grid[i][j]
			if e != 'A' {
				continue
			}
			a := grid[i-1][j-1]
			b := grid[i+1][j+1]
			if a > b {
				a, b = b, a
			}
			if a != 'M' || b != 'S' {
				continue
			}
			a = grid[i-1][j+1]
			b = grid[i+1][j-1]
			if a > b {
				a, b = b, a
			}
			if a != 'M' || b != 'S' {
				continue
			}
			count2++
		}
	}
	fmt.Println("Part 2:", count2)
}

type (
	Matcher struct {
		s      int
		target string
	}
)

func (m *Matcher) Reset() {
	m.s = 0
}

func (m *Matcher) Match(c byte) bool {
	if c == m.target[m.s] {
		m.s = (m.s + 1) % len(m.target)
		if m.s == 0 {
			return true
		}
	} else if c == m.target[0] {
		m.s = 1
	} else {
		m.s = 0
	}
	return false
}
