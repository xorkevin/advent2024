package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

	var sum uint64 = 0
	seqToPrice := map[uint32]uint32{}
	priceMap := map[uint32]uint32{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, err := strconv.ParseUint(scanner.Text(), 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		clear(priceMap)
		sum += getPriceSeq(num, seqToPrice, priceMap)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum)

	var maxNum uint32 = 0
	for _, i := range seqToPrice {
		maxNum = max(maxNum, i)
	}
	fmt.Println("Part 2:", maxNum)
}

const (
	modSeq = 19 * 19 * 19 * 19
)

func getPriceSeq(a uint64, seqToPrice, priceMap map[uint32]uint32) uint64 {
	var seq uint32 = 0
	prev := int8(a % 10)
	for i := range 2000 {
		a = calcNext(a)
		next := int8(a % 10)
		seq = ((seq * 19) % modSeq) + convertTwosComplement(next-prev)
		prev = next
		if i >= 3 {
			if _, ok := priceMap[seq]; !ok {
				priceMap[seq] = uint32(next)
			}
		}
	}
	for k, v := range priceMap {
		seqToPrice[k] += v
	}
	return a
}

func convertTwosComplement(a int8) uint32 {
	if a < 0 {
		return uint32(-a) + 9
	}
	return uint32(a)
}

const (
	modBase = 16777216
)

func calcNext(a uint64) uint64 {
	a = (a ^ (a << 6)) % modBase
	a = (a ^ (a >> 5)) % modBase
	return (a ^ (a << 11)) % modBase
}
