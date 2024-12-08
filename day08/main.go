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

	nodes := map[byte][]Pos{}
	h := 0
	w := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w = len(scanner.Bytes())
		for c, i := range scanner.Bytes() {
			if i == '.' {
				continue
			}
			nodes[i] = append(nodes[i], Pos{
				x: c,
				y: h,
			})
		}
		h++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	seen1 := map[Pos]struct{}{}
	seen2 := map[Pos]struct{}{}
	for _, freq := range nodes {
		for n, i := range freq[:len(freq)-1] {
			for _, j := range freq[n+1:] {
				delta := subPos(i, j)
				if p := addPos(i, delta); inBounds(p, w, h) {
					seen1[p] = struct{}{}
				}
				if p := subPos(j, delta); inBounds(p, w, h) {
					seen1[p] = struct{}{}
				}
				for p := i; inBounds(p, w, h); p = addPos(p, delta) {
					seen2[p] = struct{}{}
				}
				for p := j; inBounds(p, w, h); p = subPos(p, delta) {
					seen2[p] = struct{}{}
				}
			}
		}
	}
	fmt.Println("Part 1:", len(seen1))
	fmt.Println("Part 2:", len(seen2))
}

type (
	Pos struct {
		x, y int
	}
)

func addPos(a, b Pos) Pos {
	return Pos{
		x: a.x + b.x,
		y: a.y + b.y,
	}
}

func subPos(a, b Pos) Pos {
	return Pos{
		x: a.x - b.x,
		y: a.y - b.y,
	}
}

func inBounds(pos Pos, w, h int) bool {
	return pos.x >= 0 && pos.y >= 0 && pos.x < w && pos.y < h
}
