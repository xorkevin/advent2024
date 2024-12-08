pub struct BitSet {
    bits: Vec<usize>,
    size: usize,
}

impl BitSet {
    pub fn new(size: usize) -> Self {
        Self {
            bits: vec![0; (size + (usize::BITS as usize - 1)) / usize::BITS as usize],
            size: 0,
        }
    }

    pub fn len(&self) -> usize {
        self.size
    }

    pub fn reset(&mut self) {
        for i in self.bits.iter_mut() {
            *i = 0;
        }
        self.size = 0;
    }

    pub fn contains(&self, i: usize) -> bool {
        let a = i / usize::BITS as usize;
        let b = i % usize::BITS as usize;
        let mask = 1 << b;
        self.bits[a] & mask != 0
    }

    pub fn insert(&mut self, i: usize) -> bool {
        let a = i / usize::BITS as usize;
        let mask = 1 << (i % usize::BITS as usize);
        let diff = (self.bits[a] & mask) == 0;
        self.bits[a] |= mask;
        if diff {
            self.size += 1;
        }
        diff
    }

    pub fn remove(&mut self, i: usize) -> bool {
        let a = i / usize::BITS as usize;
        let mask = 1 << (i % usize::BITS as usize);
        let diff = (self.bits[a] & mask) != 0;
        self.bits[a] &= !mask;
        if diff {
            self.size -= 1;
        }
        diff
    }
}
