use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut grid = Vec::new();

    for line in reader.lines() {
        let line = line?;
        grid.push(line.into_bytes());
    }

    let grid = grid;

    let mut count1 = 0;
    let mut m1 = Matcher::new(b"XMAS");
    let mut m2 = Matcher::new(b"SAMX");
    let h = grid.len();
    let w = grid[0].len();
    for i in 0..h {
        m1.reset();
        m2.reset();
        for j in 0..w {
            let c = grid[i][j];
            if m1.read(c) {
                count1 += 1;
            }
            if m2.read(c) {
                count1 += 1;
            }
        }
    }
    for j in 0..w {
        m1.reset();
        m2.reset();
        for i in 0..h {
            let c = grid[i][j];
            if m1.read(c) {
                count1 += 1;
            }
            if m2.read(c) {
                count1 += 1;
            }
        }
    }
    for j in 0..w {
        m1.reset();
        m2.reset();
        let mut k = 0;
        while k < h && j + k < w {
            let c = grid[k][j + k];
            if m1.read(c) {
                count1 += 1;
            }
            if m2.read(c) {
                count1 += 1;
            }
            k += 1;
        }
    }
    for i in 1..h {
        m1.reset();
        m2.reset();
        let mut k = 0;
        while i + k < h && k < w {
            let c = grid[i + k][k];
            if m1.read(c) {
                count1 += 1;
            }
            if m2.read(c) {
                count1 += 1;
            }
            k += 1;
        }
    }
    for j in 0..w {
        m1.reset();
        m2.reset();
        let mut k = 0;
        while k < h && j.wrapping_sub(k) < w {
            let c = grid[k][j - k];
            if m1.read(c) {
                count1 += 1;
            }
            if m2.read(c) {
                count1 += 1;
            }
            k += 1;
        }
    }
    for i in 1..h {
        m1.reset();
        m2.reset();
        let mut k = 0;
        while i + k < h && w.wrapping_sub(k + 1) < w {
            let c = grid[i + k][w - k - 1];
            if m1.read(c) {
                count1 += 1;
            }
            if m2.read(c) {
                count1 += 1;
            }
            k += 1;
        }
    }
    println!("Part 1: {}", count1);

    let mut count2 = 0;
    for i in 1..h - 1 {
        for j in 1..w - 1 {
            let e = grid[i][j];
            if e != b'A' {
                continue;
            }
            let mut a = grid[i - 1][j - 1];
            let mut b = grid[i + 1][j + 1];
            if a > b {
                (a, b) = (b, a);
            }
            if a != b'M' || b != b'S' {
                continue;
            }
            let mut a = grid[i - 1][j + 1];
            let mut b = grid[i + 1][j - 1];
            if a > b {
                (a, b) = (b, a);
            }
            if a != b'M' || b != b'S' {
                continue;
            }
            count2 += 1;
        }
    }
    println!("Part 2: {}", count2);

    Ok(())
}

struct Matcher<const N: usize> {
    s: usize,
    target: &'static [u8; N],
}

impl<const N: usize> Matcher<N> {
    fn new(target: &'static [u8; N]) -> Self {
        Self { s: 0, target }
    }

    fn reset(&mut self) {
        self.s = 0;
    }

    fn read(&mut self, c: u8) -> bool {
        if c == self.target[self.s] {
            self.s = (self.s + 1) % N;
            if self.s == 0 {
                return true;
            }
        } else if c == self.target[0] {
            self.s = 1;
        } else {
            self.s = 0;
        }
        false
    }
}
