package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
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

	var nums1, nums2 []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f := strings.Fields(scanner.Text())
		num1, err := strconv.Atoi(f[0])
		if err != nil {
			log.Fatalln(err)
		}
		num2, err := strconv.Atoi(f[1])
		if err != nil {
			log.Fatalln(err)
		}
		nums1 = append(nums1, num1)
		nums2 = append(nums2, num2)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	slices.Sort(nums1)
	slices.Sort(nums2)

	minNum2 := nums2[0]
	maxNum2 := nums2[len(nums2)-1]
	counts := make([]int, maxNum2-minNum2+1)
	for _, i := range nums2 {
		counts[i-minNum2] += 1
	}

	s1 := 0
	s2 := 0
	for n, i := range nums1 {
		s1 += abs(i - nums2[n])
		if i < minNum2 {
			continue
		}
		s2 += i * counts[i-minNum2]
	}

	fmt.Println("Part 1:", s1)
	fmt.Println("Part 2:", s2)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
