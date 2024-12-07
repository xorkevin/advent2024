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

	count1 := 0
	count2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lhs, rhs, ok := strings.Cut(scanner.Text(), ": ")
		if !ok {
			log.Fatalln("Invalid line")
		}
		f := strings.Fields(rhs)
		nums := make([]int, 0, len(f))
		target, err := strconv.Atoi(lhs)
		if err != nil {
			log.Fatalln(err)
		}
		for _, i := range f {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatalln(err)
			}
			nums = append(nums, num)
		}
		if search(target, nums[0], nums[1:], false) {
			count1 += target
		}
		if search(target, nums[0], nums[1:], true) {
			count2 += target
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", count1)
	fmt.Println("Part 2:", count2)
}

func search(target, start int, rest []int, pt2 bool) bool {
	if len(rest) == 0 {
		return target == start
	}
	if search(target, start+rest[0], rest[1:], pt2) {
		return true
	}
	if search(target, start*rest[0], rest[1:], pt2) {
		return true
	}
	if !pt2 {
		return false
	}
	return search(target, start*magnitude(rest[0])+rest[0], rest[1:], pt2)
}

func magnitude(i int) int {
	mag := 10
	for i >= 10 {
		i /= 10
		mag *= 10
	}
	return mag
}
