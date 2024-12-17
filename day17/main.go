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

	var reg [3]int
	var program []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		lhs, rhs, ok := strings.Cut(line, ": ")
		if !ok {
			continue
		}
		if lhs == "Program" {
			f := strings.Split(rhs, ",")
			program = make([]int, 0, len(f))
			for _, i := range f {
				num, err := strconv.Atoi(i)
				if err != nil {
					log.Fatalln(err)
				}
				program = append(program, num)
			}
			continue
		}
		if strings.HasPrefix(lhs, "Register ") {
			num, err := strconv.Atoi(rhs)
			if err != nil {
				log.Fatalln(err)
			}
			switch {
			case strings.HasSuffix(lhs, "A"):
				reg[0] = num
			case strings.HasSuffix(lhs, "B"):
				reg[1] = num
			case strings.HasSuffix(lhs, "C"):
				reg[2] = num
			default:
				log.Fatalln("Invalid register")
			}
			continue
		}
		log.Fatalln("Invalid line")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	output := Exec2(reg[0])
	var sb strings.Builder
	if len(output) > 0 {
		fmt.Fprint(&sb, output[0])
		for _, i := range output[1:] {
			sb.WriteByte(',')
			fmt.Fprint(&sb, i)
		}
	}
	fmt.Println("Part 1:", sb.String())
	slices.Reverse(program)
	actual, ok := ExecRev(0, program)
	if !ok {
		log.Fatalln("Failed to invert input")
	}
	fmt.Println("Part 2:", actual)
}

func ExecRev(a int, targets []int) (int, bool) {
	if len(targets) == 0 {
		return a, true
	}
	target := targets[0]
	a <<= 3
	for b := 0; b < 8; b++ {
		if b^4^(((a|b)>>(b^1))&0b111) == target {
			if v, ok := ExecRev(a|b, targets[1:]); ok {
				return v, true
			}
		}
	}
	return 0, false
}

func ExecRevTest(output []int) int {
	var a int
	for _, i := range output {
		a <<= 3
		a &^= 0b111
		a |= i
	}
	return a
}

/*
2,4 bst A
1,1 bxl 1
7,5 cdv B
1,5 bxl 5
0,3 adv 3
4,3 bxc
5,5 out B
3,0 jnz 0
*/

/*
B = A % 8
B = B ^ 1
C = A >> B
B = B ^ 5
A = A >> 3
B = B ^ C
out B%8
jnz A 0
*/

func Exec2(a int) []int {
	var output []int
	for a != 0 {
		b := a & 0b111
		o := b ^ 4 ^ ((a >> (b ^ 1)) & 0b111)
		output = append(output, o)
		a = a >> 3
	}
	return output
}

type (
	Machine struct {
		pc      int
		reg     [3]int
		program []int
		output  []int
	}
)

func (m *Machine) Reset(a int) {
	m.pc = 0
	m.reg = [3]int{a, 0, 0}
	m.output = m.output[:0]
}

func (m *Machine) Exec() {
	for m.pc < len(m.program) {
		m.pc = m.step(m.program[m.pc], m.program[m.pc+1])
	}
}

func (m *Machine) step(op int, operand int) int {
	switch op {
	case 0:
		m.reg[0] >>= m.getComboOp(operand)
		return m.pc + 2
	case 1:
		m.reg[1] ^= operand
		return m.pc + 2
	case 2:
		m.reg[1] = m.getComboOp(operand) & 0b111
		return m.pc + 2
	case 3:
		if m.reg[0] == 0 {
			return m.pc + 2
		}
		return operand
	case 4:
		m.reg[1] ^= m.reg[2]
		return m.pc + 2
	case 5:
		m.output = append(m.output, m.getComboOp(operand)&0b111)
		return m.pc + 2
	case 6:
		m.reg[1] = m.reg[0] >> m.getComboOp(operand)
		return m.pc + 2
	case 7:
		m.reg[2] = m.reg[0] >> m.getComboOp(operand)
		return m.pc + 2
	default:
		log.Fatalln("Invalid op")
		return 0
	}
}

func (m *Machine) getComboOp(i int) int {
	switch i {
	case 0, 1, 2, 3:
		return i
	case 4:
		return m.reg[0]
	case 5:
		return m.reg[1]
	case 6:
		return m.reg[2]
	default:
		log.Fatalln("Invalid combo op")
		return 0
	}
}
