package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	isPhase1 := os.Getenv("PHASE") == "1"

	sum1 := uint64(0)
	p2Offset := Vec2{x: 10000000000000, y: 10000000000000}

	var buttonA, buttonB Vec2
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lhs, rhs, ok := strings.Cut(scanner.Text(), ": ")
		if !ok {
			continue
		}
		if strings.HasPrefix(lhs, "Button ") {
			button := parseButton(rhs)
			if strings.HasSuffix(lhs, "A") {
				buttonA = button
			} else if strings.HasSuffix(lhs, "B") {
				buttonB = button
			} else {
				log.Fatalln("Invalid line")
			}
		} else if lhs == "Prize" {
			prize := parsePrize(rhs)
			if isPhase1 {
				prize2 := addVec2(prize, p2Offset)
				fmt.Printf("%d,%d,%d,%d,%d,%d\n", prize2.x, prize2.y, buttonA.x, buttonA.y, buttonB.x, buttonB.y)
				continue
			}
			if c := search(prize, buttonA, buttonB); c < uint64(1)<<63 {
				sum1 += c
			}
		} else {
			log.Fatalln("Invalid line")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	if isPhase1 {
		return
	}

	fmt.Println("Part 1:", sum1)
}

func search(target, a, b Vec2) uint64 {
	minCost := uint64(1) << 63
	var prevV0 Vec2
	for i := uint64(0); i <= 100; i++ {
		var v0 Vec2
		if i != 0 {
			v0 = addVec2(prevV0, a)
		}
		prevV0 = v0
		if v0 == target {
			minCost = min(minCost, i*3)
			break
		}
		if gtVec2(v0, target) {
			break
		}
		prevV := v0
		for j := uint64(1); j <= 100; j++ {
			v := addVec2(prevV, b)
			prevV = v
			if v == target {
				minCost = min(minCost, i*3+j)
				break
			}
			if gtVec2(v, target) {
				break
			}
		}
	}
	return minCost
}

type (
	Vec2 struct {
		x, y int
	}
)

func addVec2(a, b Vec2) Vec2 {
	return Vec2{x: a.x + b.x, y: a.y + b.y}
}

func gtVec2(a, b Vec2) bool {
	return a.x > b.x || a.y > b.y
}

func parsePrize(s string) Vec2 {
	xs, ys, ok := strings.Cut(s, ", ")
	if !ok {
		log.Fatalln("Invalid prize", s)
	}
	if !strings.HasPrefix(xs, "X") {
		log.Fatalln("Invalid prize", s)
	}
	if !strings.HasPrefix(ys, "Y") {
		log.Fatalln("Invalid prize", s)
	}
	return Vec2{
		x: parseDelta(xs, "="),
		y: parseDelta(ys, "="),
	}
}

func parseButton(s string) Vec2 {
	xs, ys, ok := strings.Cut(s, ", ")
	if !ok {
		log.Fatalln("Invalid button", s)
	}
	if !strings.HasPrefix(xs, "X") {
		log.Fatalln("Invalid button", s)
	}
	if !strings.HasPrefix(ys, "Y") {
		log.Fatalln("Invalid button", s)
	}
	return Vec2{
		x: parseDelta(xs, "+"),
		y: parseDelta(ys, "+"),
	}
}

func parseDelta(s string, sep string) int {
	_, rhs, ok := strings.Cut(s, sep)
	if !ok {
		log.Fatalln("Invalid button delta", s)
	}
	num, err := strconv.Atoi(rhs)
	if err != nil {
		log.Fatalln(err)
	}
	return num
}
