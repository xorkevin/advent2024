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
	orders := make([][]int, 100)

	firstHalf := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if firstHalf {
			if len(scanner.Bytes()) == 0 {
				firstHalf = false
				continue
			}
			f := strings.Split(scanner.Text(), "|")
			num1, err := strconv.Atoi(f[0])
			if err != nil {
				log.Fatalln(err)
			}
			num2, err := strconv.Atoi(f[1])
			if err != nil {
				log.Fatalln(err)
			}
			// num1 must be ordered before num2
			orders[num1] = append(orders[num1], num2)
			continue
		}
		f := strings.Split(scanner.Text(), ",")
		nums := make([]int, 0, len(f))
		for _, i := range f {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatalln(err)
			}
			nums = append(nums, num)
		}

		inOrder := true
		seen := make([]bool, 100)
	outer:
		for _, i := range nums {
			for _, o := range orders[i] {
				if seen[o] {
					inOrder = false
					break outer
				}
			}
			seen[i] = true
		}
		if inOrder {
			count1 += nums[len(nums)/2]
		} else {
			slices.SortFunc(nums, func(a, b int) int {
				if a == b {
					return 0
				}
				for _, i := range orders[a] {
					if i == b {
						return -1
					}
				}
				for _, i := range orders[b] {
					if i == a {
						return 1
					}
				}
				return 0
			})
			count2 += nums[len(nums)/2]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(count1, count2)
}
