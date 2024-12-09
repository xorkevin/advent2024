package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

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

	count1 := 0
	count2 := 0
	orders := [100][100]byte{}
	seen := bitset.New(100)

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
			current := orders[num1][0] + 1
			orders[num1][current] = byte(num2)
			orders[num1][0] = current
			continue
		}
		f := strings.Split(scanner.Text(), ",")
		nums := make([]byte, 0, len(f))
		for _, i := range f {
			num, err := strconv.Atoi(i)
			if err != nil {
				log.Fatalln(err)
			}
			nums = append(nums, byte(num))
		}

		inOrder := true
		seen.Reset()
	outer:
		for _, i := range nums {
			for _, o := range orders[i][1 : orders[i][0]+1] {
				if seen.Contains(int(o)) {
					inOrder = false
					break outer
				}
			}
			seen.Insert(int(i))
		}
		if inOrder {
			count1 += int(nums[len(nums)/2])
		} else {
			slices.SortFunc(nums, func(a, b byte) int {
				if a == b {
					return 0
				}
				for _, i := range orders[a][1 : orders[a][0]+1] {
					if i == b {
						return -1
					}
				}
				for _, i := range orders[b][1 : orders[b][0]+1] {
					if i == a {
						return 1
					}
				}
				return 0
			})
			count2 += int(nums[len(nums)/2])
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", count1)
	fmt.Println("Part 2:", count2)
}
