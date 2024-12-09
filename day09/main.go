package main

import (
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	log.SetFlags(log.Lshortfile)

	file, err := os.ReadFile(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	freeListSize := len(file) / 2
	numFiles := len(file) - freeListSize
	freeList1 := make([]Segment, 0, freeListSize)
	files1 := make([]Segment, 0, numFiles)
	free := false
	id := 0
	pos := 0
	for _, i := range file {
		s := int(i - '0')
		if s > 9 {
			continue
		}
		if free {
			free = false
			freeList1 = append(freeList1, Segment{id: 0, pos: pos, size: s})
		} else {
			free = true
			files1 = append(files1, Segment{id: id, pos: pos, size: s})
			id++
		}
		pos += s
	}

	freeList2 := make([]Segment, len(freeList1))
	copy(freeList2, freeList1)
	files2 := make([]Segment, len(files1))
	copy(files2, files1)

	nextFiles1 := make([]Segment, 0, numFiles)
	for len(freeList1) > 0 && freeList1[0].pos < files1[len(files1)-1].pos {
		freeNode := freeList1[0]
		f := files1[len(files1)-1]
		if f.size == freeNode.size {
			f.pos = freeNode.pos
			nextFiles1 = append(nextFiles1, f)
			freeList1 = freeList1[1:]
			files1 = files1[:len(files1)-1]
		} else if f.size < freeNode.size {
			f.pos = freeNode.pos
			nextFiles1 = append(nextFiles1, f)
			freeNode.pos += f.size
			freeNode.size -= f.size
			freeList1[0] = freeNode
			files1 = files1[:len(files1)-1]
		} else {
			freeNode.id = f.id
			nextFiles1 = append(nextFiles1, freeNode)
			freeList1 = freeList1[1:]
			files1[len(files1)-1].size -= freeNode.size
		}
	}

	sum1 := 0
	for _, i := range files1 {
		sum1 += i.id * (i.size * (i.pos + i.pos + (i.size - 1)) / 2)
	}
	for _, i := range nextFiles1 {
		sum1 += i.id * (i.size * (i.pos + i.pos + (i.size - 1)) / 2)
	}
	fmt.Println("Part 1:", sum1)

	for i := len(files2) - 1; i >= 0; i-- {
		f := files2[i]
		for j, freeNode := range freeList2 {
			if freeNode.pos >= f.pos {
				break
			}
			if f.size <= freeNode.size {
				f.pos = freeNode.pos
				files2[i] = f
				freeNode.pos += f.size
				freeNode.size -= f.size
				freeList2[j] = freeNode
				break
			}
		}
	}

	sum2 := 0
	for _, i := range files2 {
		sum2 += i.id * (i.size * (i.pos + i.pos + (i.size - 1)) / 2)
	}
	fmt.Println("Part 2:", sum2)
}

type (
	Segment struct {
		id   int
		pos  int
		size int
	}
)
