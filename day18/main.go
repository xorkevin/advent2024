package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/xorkevin/advent2024/graph"
)

const (
	puzzleInput = "input.txt"
	// dim         = 7
	dim = 71
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

	var grid [dim][dim]byte
	for i, row := range grid {
		for j := range row {
			grid[i][j] = '.'
		}
	}

	count := 0
	var rest []Vec2
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lhs, rhs, ok := strings.Cut(scanner.Text(), ",")
		if !ok {
			log.Fatalln("Invalid line")
		}
		x, err := strconv.Atoi(lhs)
		if err != nil {
			log.Fatalln(err)
		}
		y, err := strconv.Atoi(rhs)
		if err != nil {
			log.Fatalln(err)
		}
		if count < 1024 {
			grid[y][x] = '#'
			count++
			continue
		}
		rest = append(rest, Vec2{x: x, y: y})
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	start := Vec2{x: 0, y: 0}
	end := Vec2{x: dim - 1, y: dim - 1}
	fmt.Println("Part 1:", graph.Search(graph.Edge[Vec2]{
		V: start,
		C: 0,
		H: manhattanDistance(start, end),
	}, func(v Vec2) bool {
		return v == end
	}, func(v Vec2) []graph.Edge[Vec2] {
		var n []graph.Edge[Vec2]
		for _, i := range neighborDeltas {
			k := addPos(v, i)
			if isOutBounds(k) {
				continue
			}
			if grid[k.y][k.x] == '#' {
				continue
			}
			n = append(n, graph.Edge[Vec2]{
				V: k,
				C: 1,
				H: manhattanDistance(k, end),
			})
		}
		return n
	}))
	for _, i := range rest {
		grid[i.y][i.x] = '#'
		if graph.Search(graph.Edge[Vec2]{
			V: start,
			C: 0,
			H: manhattanDistance(start, end),
		}, func(v Vec2) bool {
			return v == end
		}, func(v Vec2) []graph.Edge[Vec2] {
			var n []graph.Edge[Vec2]
			for _, i := range neighborDeltas {
				k := addPos(v, i)
				if isOutBounds(k) {
					continue
				}
				if grid[k.y][k.x] == '#' {
					continue
				}
				n = append(n, graph.Edge[Vec2]{
					V: k,
					C: 1,
					H: manhattanDistance(k, end),
				})
			}
			return n
		}) < 0 {
			fmt.Printf("Part 2: %d,%d\n", i.x, i.y)
			return
		}
	}
}

func isOutBounds(v Vec2) bool {
	return v.x < 0 || v.y < 0 || v.x >= dim || v.y >= dim
}

var neighborDeltas = []Vec2{
	{x: 0, y: -1},
	{x: 1, y: 0},
	{x: 0, y: 1},
	{x: -1, y: 0},
}

type (
	Vec2 struct {
		x, y int
	}
)

func addPos(a, b Vec2) Vec2 {
	return Vec2{x: a.x + b.x, y: a.y + b.y}
}

func manhattanDistance(a, b Vec2) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
