package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
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

	count := 0
	count2 := 0

	isFirst := true
	var patterns []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			isFirst = false
			continue
		}
		if isFirst {
			patterns = strings.Split(scanner.Text(), ", ")
			continue
		}
		c := search(scanner.Bytes(), patterns)
		if c > 0 {
			count++
		}
		count2 += c
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", count)
	fmt.Println("Part 2:", count2)
}

func search(target []byte, patterns []string) int {
	l := len(target)
	dp := make([]int, l+1)
	dp[0] = 1
	for idx := 1; idx <= l; idx++ {
		for _, i := range patterns {
			if !bytes.HasSuffix(target[:idx], []byte(i)) {
				continue
			}
			dp[idx] += dp[idx-len(i)]
		}
	}
	return dp[l]
}
