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

	var codes [][]byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		code := make([]byte, 4)
		row := scanner.Bytes()
		for i := range code {
			code[i] = row[i] - '0'
		}
		code[3] = 10
		codes = append(codes, code)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	sum1 := 0
	sum2 := 0
	var memo1 [26]map[string]int
	var memo2 [26]map[string]int
	for i := range memo1 {
		memo1[i] = map[string]int{}
	}
	for i := range memo2 {
		memo2[i] = map[string]int{}
	}
	for _, i := range codes {
		num := 0
		for _, j := range i[:3] {
			num *= 10
			num += int(j)
		}
		sum1 += num * search(i, 2, true, memo1)
		sum2 += num * search(i, 25, true, memo2)
	}
	fmt.Println("Part 1:", sum1)
	fmt.Println("Part 2:", sum2)
}

func search(code []byte, depth int, top bool, memo [26]map[string]int) int {
	if depth < 0 {
		return len(code)
	}

	if v, ok := memo[depth][string(code)]; ok {
		return v
	}

	var k byte = 4
	if top {
		k = 10
	}
	sum := 0
	for _, i := range code {
		switch i {
		case '^':
			i = 0
		case '>':
			i = 1
		case 'v':
			i = 2
		case '<':
			i = 3
		case 'A':
			i = 4
		}
		coord := dirpadCoord
		if top {
			coord = keypadCoord
		}
		minCost := 1 << 62
		children := findDirs(k, i, coord, top)
		for _, j := range children {
			minCost = min(minCost, search([]byte(j), depth-1, false, memo))
		}
		sum += minCost
		k = i
	}
	memo[depth][string(code)] = sum
	return sum
}

func findDirs(from, to byte, coord []Vec2, mode bool) []string {
	var arr [2]string
	n := arr[:0]
	fromPos := coord[from]
	delta := subPos(coord[to], fromPos)
	horizontal := ""
	vertical := ""
	if delta.x < 0 {
		for range abs(delta.x) {
			horizontal += "<"
		}
	} else if delta.x > 0 {
		for range abs(delta.x) {
			horizontal += ">"
		}
	}
	if delta.y < 0 {
		for range abs(delta.y) {
			vertical += "^"
		}
	} else if delta.y > 0 {
		for range abs(delta.y) {
			vertical += "v"
		}
	}
	if mode {
		if addPos(fromPos, Vec2{x: 0, y: delta.y}) != (Vec2{x: 0, y: 3}) {
			n = append(n, vertical+horizontal+"A")
		}
		if addPos(fromPos, Vec2{x: delta.x, y: 0}) != (Vec2{x: 0, y: 3}) {
			n = append(n, horizontal+vertical+"A")
		}
	} else {
		if addPos(fromPos, Vec2{x: 0, y: delta.y}) != (Vec2{x: 0, y: 0}) {
			n = append(n, vertical+horizontal+"A")
		}
		if addPos(fromPos, Vec2{x: delta.x, y: 0}) != (Vec2{x: 0, y: 0}) {
			n = append(n, horizontal+vertical+"A")
		}
	}
	return n
}

var (
	dirpadCoord = []Vec2{
		{x: 1, y: 0},
		{x: 2, y: 1},
		{x: 1, y: 1},
		{x: 0, y: 1},
		{x: 2, y: 0},
	}

	keypadCoord = []Vec2{
		{x: 1, y: 3},
		{x: 0, y: 2},
		{x: 1, y: 2},
		{x: 2, y: 2},
		{x: 0, y: 1},
		{x: 1, y: 1},
		{x: 2, y: 1},
		{x: 0, y: 0},
		{x: 1, y: 0},
		{x: 2, y: 0},
		{x: 2, y: 3},
	}
)

type (
	Vec2 struct {
		x, y int
	}
)

func addPos(a, b Vec2) Vec2 {
	return Vec2{x: a.x + b.x, y: a.y + b.y}
}

func subPos(a, b Vec2) Vec2 {
	return Vec2{x: a.x - b.x, y: a.y - b.y}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
