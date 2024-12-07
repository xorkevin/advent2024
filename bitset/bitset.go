package bitset

type (
	BitSet struct {
		bits []uint64
		size int
	}
)

func New(size int) *BitSet {
	return &BitSet{
		bits: make([]uint64, (size+63)/64),
		size: 0,
	}
}

func (s *BitSet) Reset() {
	for i := range s.bits {
		s.bits[i] = 0
	}
	s.size = 0
}

func (s *BitSet) Size() int {
	return s.size
}

func (s *BitSet) Contains(i int) bool {
	a := i / 64
	mask := uint64(1) << (i % 64)
	return (s.bits[a] & mask) != 0
}

func (s *BitSet) Set(i int, b bool) bool {
	a := i / 64
	mask := uint64(1) << (i % 64)
	if b {
		diff := (s.bits[a] & mask) == 0
		s.bits[a] |= mask
		if diff {
			s.size++
		}
		return diff
	} else {
		diff := (s.bits[a] & mask) != 0
		s.bits[a] &^= mask
		if diff {
			s.size--
		}
		return diff
	}
}

func (s *BitSet) Insert(i int) bool {
	a := i / 64
	mask := uint64(1) << (i % 64)
	diff := (s.bits[a] & mask) == 0
	s.bits[a] |= mask
	if diff {
		s.size++
	}
	return diff
}

func (s *BitSet) Remove(i int) bool {
	a := i / 64
	mask := uint64(1) << (i % 64)
	diff := (s.bits[a] & mask) != 0
	s.bits[a] &^= mask
	if diff {
		s.size--
	}
	return diff
}
