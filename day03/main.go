package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

var mulInstrRegex = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do(?:n't)?\(\)`)

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

	sum1 := 0
	sum2 := 0
	enabled := true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		found := mulInstrRegex.FindAllStringSubmatch(scanner.Text(), -1)
		for _, i := range found {
			if strings.HasPrefix(i[0], "don") {
				enabled = false
				continue
			}
			if strings.HasPrefix(i[0], "d") {
				enabled = true
				continue
			}
			num1, err := strconv.Atoi(i[1])
			if err != nil {
				log.Fatalln(err)
			}
			num2, err := strconv.Atoi(i[2])
			if err != nil {
				log.Fatalln(err)
			}
			product := num1 * num2
			sum1 += product
			if enabled {
				sum2 += product
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum1)
	fmt.Println("Part 2:", sum2)
}
