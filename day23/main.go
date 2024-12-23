package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
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

	edges := map[string][]string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lhs, rhs, ok := strings.Cut(scanner.Text(), "-")
		if !ok {
			log.Fatalln("Invalid line")
		}
		if rhs < lhs {
			lhs, rhs = rhs, lhs
		}
		edges[lhs] = append(edges[lhs], rhs)
		slices.Sort(edges[lhs])
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", getCliques3(edges))
	maxSize := 0
	var clique string
	for k := range edges {
		c, size := getLargestClique(edges, nil, k)
		if size > maxSize {
			maxSize = size
			clique = c
		}
	}
	fmt.Println("Part 2:", clique)
}

func getLargestClique(edges map[string][]string, s []string, a string) (string, int) {
	s1 := append(s, a)

	maxSize := 0
	var clique string

loop:
	for _, i := range edges[a] {
		for _, j := range s {
			p := j
			q := i
			if q < p {
				p, q = q, p
			}
			if !slices.Contains(edges[p], q) {
				continue loop
			}
		}
		c, size := getLargestClique(edges, s1, i)
		if size > maxSize {
			maxSize = size
			clique = c
		}
	}
	if maxSize == 0 {
		return strings.Join(s1, ","), len(s1)
	}
	return clique, maxSize
}

func getCliques3(edges map[string][]string) int {
	m := map[string]struct{}{}
	for k, v := range edges {
		for n, i := range v[:len(v)-1] {
			for _, j := range v[n+1:] {
				if !strings.HasPrefix(k, "t") && !strings.HasPrefix(i, "t") && !strings.HasPrefix(j, "t") {
					continue
				}
				if slices.Contains(edges[i], j) {
					m[k+","+i+","+j] = struct{}{}
				}
			}
		}
	}
	return len(m)
}
