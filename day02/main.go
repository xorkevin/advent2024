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
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	count1 := 0
	count2 := 0

	var row []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f := strings.Fields(scanner.Text())
		lf := len(f)
		if len(row) < lf {
			row = make([]int, len(f))
		}
		r := row[:lf]
		for n, i := range f {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatalln(err)
			}
			r[n] = num
		}
		if isSafe(r, -1) {
			count1++
			count2++
		} else {
			for i := 0; i < lf; i++ {
				if isSafe(r, i) {
					count2++
					break
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", count1)
	fmt.Println("Part 2:", count2)
}

func isSafe(row []int, exclude int) bool {
	first := true
	second := true
	prev := 0
	inc := true
	for n, num := range row {
		if n == exclude {
			continue
		}
		if first {
			first = false
			prev = num
			continue
		}
		delta := num - prev
		prev = num
		if a := abs(delta); a < 1 || a > 3 {
			return false
		}
		pos := delta > 0
		if second {
			second = false
			inc = pos
			continue
		}
		if pos != inc {
			return false
		}
	}
	return true
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
