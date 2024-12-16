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
	needsStart := true
	var start Vec2

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if needsStart {
			x := bytes.IndexByte(scanner.Bytes(), 'S')
			if x >= 0 {
				needsStart = false
				start = Vec2{x: x, y: len(grid)}
			}
		}
		grid = append(grid, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	bestCost, bestPathSetSize := search(grid, State{
		pos:   start,
		dir:   1,
		cost:  0,
		depth: 0,
	})
	fmt.Println("Part 1:", bestCost)
	fmt.Println("Part 2:", bestPathSetSize+1)
}

func search(grid [][]byte, startState State) (int, int) {
	w := len(grid[0])
	stack := []State{startState}
	var openPath []Vec2
	seen := bitset.New(len(grid) * w)
	bestPathSet := bitset.New(len(grid) * w)
	bestCost := len(grid) * w * 1000
	bestCosts := make([]int, len(grid)*w*4)
	for i := range bestCosts {
		bestCosts[i] = bestCost
	}
loop:
	for len(stack) != 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if len(openPath) >= top.depth {
			for _, i := range openPath[top.depth:] {
				seen.Remove(getPosID(i, w))
			}
			openPath = openPath[:top.depth]
		}
		switch grid[top.pos.y][top.pos.x] {
		case 'E':
			if top.cost <= bestCost {
				if top.cost < bestCost {
					bestCost = top.cost
					bestPathSet.Reset()
				}
				for _, i := range openPath {
					bestPathSet.Insert(getPosID(i, w))
				}
			}
			continue loop
		case '.', 'S':
		default:
			continue loop
		}
		stateID := getStateID(top.pos, w, top.dir)
		if bestCosts[stateID] < top.cost {
			continue
		}
		if !seen.Insert(getPosID(top.pos, w)) {
			continue
		}
		openPath = append(openPath, top.pos)
		bestCosts[stateID] = top.cost
		stack = append(stack, State{
			pos:   addPos(top.pos, dirToDelta(top.dir)),
			dir:   top.dir,
			cost:  top.cost + 1,
			depth: top.depth + 1,
		}, State{
			pos:   addPos(top.pos, dirToDelta((top.dir+1)%4)),
			dir:   (top.dir + 1) % 4,
			cost:  top.cost + 1001,
			depth: top.depth + 1,
		}, State{
			pos:   addPos(top.pos, dirToDelta((top.dir+3)%4)),
			dir:   (top.dir + 3) % 4,
			cost:  top.cost + 1001,
			depth: top.depth + 1,
		})
	}
	return bestCost, bestPathSet.Size()
}

type (
	State struct {
		pos   Vec2
		dir   int
		cost  int
		depth int
	}

	Vec2 struct {
		x, y int
	}
)

func addPos(a, b Vec2) Vec2 {
	return Vec2{x: a.x + b.x, y: a.y + b.y}
}

func dirToDelta(dir int) Vec2 {
	switch dir {
	case 0:
		return Vec2{x: 0, y: -1}
	case 1:
		return Vec2{x: 1, y: 0}
	case 2:
		return Vec2{x: 0, y: 1}
	case 3:
		return Vec2{x: -1, y: 0}
	default:
		log.Fatalln("Invalid dir")
		return Vec2{}
	}
}

func getStateID(pos Vec2, w int, dir int) int {
	return getPosID(pos, w)*4 + dir
}

func getPosID(pos Vec2, w int) int {
	return pos.y*w + pos.x
}
