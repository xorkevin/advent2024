package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	var locks []uint32
	var keys []uint32

	var buf [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			if k, isLock := parseBlock(buf); isLock {
				locks = append(locks, k)
			} else {
				keys = append(keys, k)
			}
			buf = buf[:0]
			continue
		}
		buf = append(buf, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	if k, isLock := parseBlock(buf); isLock {
		locks = append(locks, k)
	} else {
		keys = append(keys, k)
	}

	count := 0
	for _, i := range locks {
		for _, j := range keys {
			if fits(i, j) {
				count++
			}
		}
	}
	fmt.Println("Part 1:", count)
}

func fits(a, b uint32) bool {
	k := a + b
	for range 5 {
		if (k & 0xf) > 5 {
			return false
		}
		k >>= 4
	}
	return true
}

func parseBlock(block [][]byte) (uint32, bool) {
	var k uint32
	isLock := block[0][0] == '#'
	for i := 0; i < 5; i++ {
		var j uint32 = 0
		for ; j < 5; j++ {
			if isLock {
				if block[1+j][i] != '#' {
					break
				}
			} else {
				if block[5-j][i] != '#' {
					break
				}
			}
		}
		k <<= 4
		k += j
	}
	return k, isLock
}
