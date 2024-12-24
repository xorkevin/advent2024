package main

import (
	"bufio"
	"fmt"
	"log"
	"math/bits"
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

	var inpX, inpY uint64
	first := true

	var allWires []WireOp

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			first = false
			continue
		}
		if first {
			lhs, rhs, ok := strings.Cut(scanner.Text(), ": ")
			if !ok {
				log.Fatalln("Invalid line")
			}
			var b uint64
			if rhs == "1" {
				b = 1
			}
			num, err := strconv.Atoi(lhs[1:])
			if err != nil {
				log.Fatalln(err)
			}
			b <<= num
			if strings.HasPrefix(lhs, "x") {
				inpX |= b
			} else {
				inpY |= b
			}
			continue
		}
		lhs, rhs, ok := strings.Cut(scanner.Text(), " -> ")
		if !ok {
			log.Fatalln("Invalid line")
		}
		f := strings.Fields(lhs)
		if len(f) != 3 {
			log.Fatalln("Invalid line")
		}
		var op WireOp
		if strings.HasPrefix(f[0], "x") || strings.HasPrefix(f[0], "y") {
			num, err := strconv.Atoi(f[0][1:])
			if err != nil {
				log.Fatalln("Invalid line")
			}
			op.inp1 = num
			op.arg1 = "x"
			if strings.HasPrefix(f[0], "y") {
				op.arg1 = "y"
			}
		} else {
			op.arg1 = f[0]
		}
		if strings.HasPrefix(f[2], "x") || strings.HasPrefix(f[2], "y") {
			num, err := strconv.Atoi(f[2][1:])
			if err != nil {
				log.Fatalln("Invalid line")
			}
			op.inp2 = num
			op.arg2 = "x"
			if strings.HasPrefix(f[2], "y") {
				op.arg2 = "y"
			}
		} else {
			op.arg2 = f[2]
		}
		switch f[1] {
		case "AND":
			op.op = 0
		case "OR":
			op.op = 1
		case "XOR":
			op.op = 2
		default:
			log.Fatalln("Invalid op")
		}
		if strings.HasPrefix(rhs, "z") {
			num, err := strconv.Atoi(rhs[1:])
			if err != nil {
				log.Fatalln("Invalid line")
			}
			op.outZ = num
			op.out = rhs
		} else {
			op.out = rhs
		}
		allWires = append(allWires, op)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	wires := map[string]byte{}
	fmt.Println("Part 1:", evalWires(inpX, inpY, allWires, wires, nil))

	swaps := map[string]string{
		"z08": "thm",
		"thm": "z08",
		"wss": "wrm",
		"wrm": "wss",
		"z22": "hwq",
		"hwq": "z22",
		"gbs": "z29",
		"z29": "gbs",
	}
	sampleAddDiff(allWires, wires, swaps)
	fmt.Println("Part 2:", strings.Join([]string{
		"gbs",
		"hwq",
		"thm",
		"wrm",
		"wss",
		"z08",
		"z22",
		"z29",
	}, ","))
}

func findBestSwap(allWires []WireOp, wires map[string]byte, swaps map[string]string) {
	minDiff := 64
	for i := 0; i < len(allWires)-1; i++ {
		a := allWires[i].out
		if _, ok := swaps[a]; ok {
			continue
		}
		for j := i + 1; j < len(allWires); j++ {
			b := allWires[j].out
			if _, ok := swaps[b]; ok {
				continue
			}
			swaps[a] = b
			swaps[b] = a
			diff := sampleAddDiff(allWires, wires, swaps)
			count := bits.OnesCount64(diff)
			if count <= minDiff {
				minDiff = count
				fmt.Printf("minDiff %d %s %s %046b\n", minDiff, a, b, diff)
			}
			delete(swaps, a)
			delete(swaps, b)
		}
	}
}

const (
	bitmask45     = 0x1fffffffffff
	everyOtherBit = 0x155555555555
	triplets      = 0b100100100100100100100100100100100100100100100
	quads         = 0x111111111111
	penta         = 0b100001000010000100001000010000100001000010000
	hepta         = 0b100000010000001000000100000010000001000000100
)

type (
	Vec2 struct {
		x, y uint64
	}
)

var testCases = []Vec2{
	{x: 0, y: 0},
	{x: 0, y: everyOtherBit},
	{x: 0, y: everyOtherBit << 1},
	{x: 0, y: bitmask45},
	{x: 0, y: triplets},
	{x: 0, y: quads},
	{x: 0, y: penta},
	{x: 0, y: hepta},
	{x: everyOtherBit, y: 0},
	{x: everyOtherBit, y: everyOtherBit},
	{x: everyOtherBit, y: everyOtherBit << 1},
	{x: everyOtherBit, y: bitmask45},
	{x: everyOtherBit, y: triplets},
	{x: everyOtherBit, y: quads},
	{x: everyOtherBit, y: penta},
	{x: everyOtherBit, y: hepta},
	{x: everyOtherBit << 1, y: 0},
	{x: everyOtherBit << 1, y: everyOtherBit},
	{x: everyOtherBit << 1, y: everyOtherBit << 1},
	{x: everyOtherBit << 1, y: bitmask45},
	{x: everyOtherBit << 1, y: triplets},
	{x: everyOtherBit << 1, y: quads},
	{x: everyOtherBit << 1, y: penta},
	{x: everyOtherBit << 1, y: hepta},
	{x: bitmask45, y: 0},
	{x: bitmask45, y: everyOtherBit},
	{x: bitmask45, y: everyOtherBit << 1},
	{x: bitmask45, y: bitmask45},
	{x: bitmask45, y: triplets},
	{x: bitmask45, y: quads},
	{x: bitmask45, y: penta},
	{x: bitmask45, y: hepta},
	{x: triplets, y: 0},
	{x: triplets, y: everyOtherBit},
	{x: triplets, y: everyOtherBit << 1},
	{x: triplets, y: bitmask45},
	{x: triplets, y: triplets},
	{x: triplets, y: quads},
	{x: triplets, y: penta},
	{x: triplets, y: hepta},
	{x: quads, y: 0},
	{x: quads, y: everyOtherBit},
	{x: quads, y: everyOtherBit << 1},
	{x: quads, y: bitmask45},
	{x: quads, y: triplets},
	{x: quads, y: quads},
	{x: quads, y: penta},
	{x: quads, y: hepta},
	{x: penta, y: 0},
	{x: penta, y: everyOtherBit},
	{x: penta, y: everyOtherBit << 1},
	{x: penta, y: bitmask45},
	{x: penta, y: triplets},
	{x: penta, y: quads},
	{x: penta, y: penta},
	{x: penta, y: hepta},
	{x: hepta, y: 0},
	{x: hepta, y: everyOtherBit},
	{x: hepta, y: everyOtherBit << 1},
	{x: hepta, y: bitmask45},
	{x: hepta, y: triplets},
	{x: hepta, y: quads},
	{x: hepta, y: penta},
	{x: hepta, y: hepta},
}

func sampleAddDiff(allWires []WireOp, wires map[string]byte, swaps map[string]string) uint64 {
	var diff uint64
	for _, i := range testCases {
		a := i.x & bitmask45
		b := i.y & bitmask45
		diff |= (evalWires(a, b, allWires, wires, swaps) ^ (a + b))
	}
	return diff
}

type (
	WireOp struct {
		inp1 int
		arg1 string
		op   byte
		inp2 int
		arg2 string
		out  string
		outZ int
	}
)

func evalWires(inpX, inpY uint64, allWires []WireOp, wires map[string]byte, swaps map[string]string) uint64 {
	clear(wires)

	for {
		processedValue := false
		for _, i := range allWires {
			o := i.out
			if v, ok := swaps[o]; ok {
				o = v
			}
			if _, ok := wires[o]; ok {
				continue
			}
			var a, b byte
			if i.arg1 == "x" || i.arg1 == "y" {
				v := inpX
				if i.arg1 == "y" {
					v = inpY
				}
				a = byte((v >> i.inp1) & 1)
			} else {
				var ok bool
				a, ok = wires[i.arg1]
				if !ok {
					continue
				}
			}
			if i.arg2 == "x" || i.arg2 == "y" {
				v := inpX
				if i.arg2 == "y" {
					v = inpY
				}
				b = byte((v >> i.inp2) & 1)
			} else {
				var ok bool
				b, ok = wires[i.arg2]
				if !ok {
					continue
				}
			}
			var out byte
			switch i.op {
			case 0:
				out = a & b
			case 1:
				out = a | b
			case 2:
				out = a ^ b
			default:
				log.Fatalln("Invalid op", i.op)
			}
			wires[o] = out
			processedValue = true
		}
		if !processedValue {
			break
		}
	}
	var val uint64
	for _, i := range allWires {
		if strings.HasPrefix(i.out, "z") {
			val |= uint64(wires[i.out]) << i.outZ
		}
	}
	return val
}
