package crt

import (
	"math/bits"
)

func CRT(a1, m1, a2, m2 int) (int, int, bool) {
	g, p, q := extGCD(m1, m2)
	if a1%g != a2%g {
		return 0, 0, false
	}
	m1g := m1 / g
	m2g := m2 / g
	lcm := m1g * m2
	// a1 * m2/g * q + a2 * m1/g * p (mod lcm)
	x := (mulmod(mulmod(a1, m2g, lcm), q, lcm) + mulmod(mulmod(a2, m1g, lcm), p, lcm)) % lcm
	if x < 0 {
		x += lcm
	}
	return x, lcm, true
}

func extGCD(a, b int) (int, int, int) {
	x2 := 1
	x1 := 0
	y2 := 0
	y1 := 1
	// a should be larger than b
	flip := false
	if a < b {
		a, b = b, a
		flip = true
	}
	for b > 0 {
		q := a / b
		a, b = b, a%b
		x2, x1 = x1, x2-q*x1
		y2, y1 = y1, y2-q*y1
	}
	if flip {
		x2, y2 = y2, x2
	}
	return a, x2, y2
}

func mulmod(a, b, m int) int {
	sign := 1
	if a < 0 {
		a = -a
		sign *= -1
	}
	if b < 0 {
		b = -b
		sign *= -1
	}
	a = a % m
	b = b % m
	hi, lo := bits.Mul(uint(a), uint(b))
	return sign * int(bits.Rem(hi, lo, uint(m)))
}
