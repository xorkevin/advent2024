package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	puzzleInput = `0 37551 469 63 1 791606 2065 9983586`
)

func main() {
	log.SetFlags(log.Lshortfile)

	f := strings.Fields(puzzleInput)
	s1 := 0
	s2 := 0
	for _, i := range f {
		num, err := strconv.Atoi(i)
		if err != nil {
			log.Fatalln(err)
		}
		s1 += nextStep(num, 25)
		s2 += nextStep(num, 75)
	}
	fmt.Println("Part 1:", s1)
	fmt.Println("Part 2:", s2)
}

type (
	Vec2 struct {
		x, y int
	}
)

var memo = map[Vec2]int{}

func nextStep(num int, depth int) int {
	if depth == 0 {
		return 1
	}

	key := Vec2{
		x: num,
		y: depth,
	}
	if res, ok := memo[key]; ok {
		return res
	}

	var res int
	if num == 0 {
		res = nextStep(1, depth-1)
	} else if a, b, ok := split(num); ok {
		res = nextStep(a, depth-1)
		res += nextStep(b, depth-1)
	} else {
		res = nextStep(2024*num, depth-1)
	}
	memo[key] = res
	return res
}

func split(n int) (int, int, bool) {
	s := strconv.Itoa(n)
	if len(s)%2 != 0 {
		return 0, 0, false
	}
	t := len(s) / 2
	a := 1
	for i := 0; i < t; i++ {
		a *= 10
	}
	return n / a, n % a, true
}
