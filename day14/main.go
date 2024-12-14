package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/xorkevin/advent2024/crt"
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

	const (
		w = 101
		h = 103
	)

	var robots []Robot
	nQuad := [4]int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		p, v := parseRobot(scanner.Text())
		idx := getQuadrant(getPos(p, v, 100, w, h), w, h)
		if idx >= 0 {
			nQuad[idx]++
		}
		robots = append(robots, Robot{p: p, v: v})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", nQuad[0]*nQuad[1]*nQuad[2]*nQuad[3])

	// short dx: multiple of 101*t + 79
	// short dy: multiple of 103*t + 22
	t2, _, ok := crt.CRT(79, 101, 22, 103)
	if !ok {
		log.Fatalln("Invalid crt")
	}
	fmt.Println("Part 2:", t2)

	if os.Getenv("PHASE") == "1" {
		t := 0
		reader := bufio.NewReader(os.Stdin)
		for {
			printGrid(w, h, robots, t)
			s, err := reader.ReadString('\n')
			if err == nil {
				num, err := strconv.Atoi(strings.TrimRight(s, "\n"))
				if err == nil {
					t = num
					continue
				}
			}
			t++
		}
	}
}

func printGrid(w, h int, robots []Robot, t int) {
	grid := make([][]byte, h)
	for n := range grid {
		row := make([]byte, w)
		for m := range row {
			row[m] = '.'
		}
		grid[n] = row
	}
	for _, i := range robots {
		p := getPos(i.p, i.v, t, w, h)
		grid[p.y][p.x] = '#'
	}
	fmt.Println(t, "s")
	for _, i := range grid {
		fmt.Println(string(i))
	}
}

type (
	Robot struct {
		p, v Vec2
	}
)

func getQuadrant(v Vec2, w, h int) int {
	w2 := w / 2
	h2 := h / 2
	if v.x == w2 || v.y == h2 {
		return -1
	}
	q := 0
	if v.x > w2 {
		q |= 1
	}
	if v.y > h2 {
		q |= 2
	}
	return q
}

func getPos(p, v Vec2, t int, w, h int) Vec2 {
	v.x = (v.x + w) % w
	v.y = (v.y + h) % h
	p.x = (p.x + v.x*t) % w
	p.y = (p.y + v.y*t) % h
	return p
}

type (
	Vec2 struct {
		x, y int
	}
)

func parseRobot(s string) (Vec2, Vec2) {
	lhs, rhs, ok := strings.Cut(s, " ")
	if !ok {
		log.Fatalln("Invalid line", s)
	}
	return parseVec2(lhs, "p="), parseVec2(rhs, "v=")
}

func parseVec2(s string, prefix string) Vec2 {
	s, ok := strings.CutPrefix(s, prefix)
	if !ok {
		log.Fatalln("Invalid vec2")
	}
	lhs, rhs, ok := strings.Cut(s, ",")
	if !ok {
		log.Fatalln("Invalid vec2")
	}
	x, err := strconv.Atoi(lhs)
	if err != nil {
		log.Fatalln(err)
	}
	y, err := strconv.Atoi(rhs)
	if err != nil {
		log.Fatalln(err)
	}
	return Vec2{x: x, y: y}
}
