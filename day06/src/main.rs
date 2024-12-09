use advent2024::bitset::BitSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut grid = Vec::new();
    let mut need_start = true;
    let mut start = Pos { x: 0, y: 0 };
    for (r, line) in reader.lines().enumerate() {
        let line = line?;
        let row = line.into_bytes();
        if need_start {
            if let Some(c) = row.iter().position(|&v| v == b'^') {
                need_start = false;
                start = Pos { x: c, y: r };
            }
        }
        grid.push(row);
    }
    let start = start;
    let grid = grid;

    let visited = sim(&grid, start);
    println!("Part 1: {}", visited.len() + 1);

    let mut count2 = 0;
    let w = grid[0].len();
    let h = grid.len();
    let mut seen = BitSet::new(w * h * 4);
    for i in visited {
        if sim_loop(&grid, start, i, &mut seen) {
            count2 += 1;
        }
    }
    println!("Part 2: {}", count2);

    Ok(())
}

fn sim_loop(grid: &Vec<Vec<u8>>, mut pos: Pos, obstruct: Pos, seen: &mut BitSet) -> bool {
    let bounds = Pos {
        x: grid[0].len(),
        y: grid.len(),
    };

    let mut delta = Delta::N;
    seen.reset();
    loop {
        let next = pos.add(delta);
        if out_bounds(next, bounds) {
            return false;
        }
        if next == obstruct || grid[next.y][next.x] == b'#' {
            delta = delta.turn();
            if !seen.insert(get_state_id(pos, bounds.x, delta)) {
                return true;
            }
            continue;
        }
        pos = next;
    }
}

fn sim(grid: &Vec<Vec<u8>>, mut pos: Pos) -> Vec<Pos> {
    let bounds = Pos {
        x: grid[0].len(),
        y: grid.len(),
    };

    let mut p = Vec::new();
    let mut delta = Delta::N;
    let mut seen = BitSet::new(bounds.x * bounds.y);
    seen.insert(get_pos_id(pos, bounds.x));
    loop {
        let next = pos.add(delta);
        if out_bounds(next, bounds) {
            break;
        }
        if grid[next.y][next.x] == b'#' {
            delta = delta.turn();
            continue;
        }
        pos = next;
        if seen.insert(get_pos_id(pos, bounds.x)) {
            p.push(pos);
        }
    }
    p
}

#[derive(Clone, Copy, PartialEq, Eq)]
struct Pos {
    x: usize,
    y: usize,
}

impl Pos {
    fn add(&self, delta: Delta) -> Self {
        match delta {
            Delta::N => Pos {
                x: self.x,
                y: self.y.wrapping_sub(1),
            },
            Delta::E => Pos {
                x: self.x + 1,
                y: self.y,
            },
            Delta::S => Pos {
                x: self.x,
                y: self.y + 1,
            },
            Delta::W => Pos {
                x: self.x.wrapping_sub(1),
                y: self.y,
            },
        }
    }
}

#[derive(Clone, Copy)]
enum Delta {
    N,
    E,
    S,
    W,
}

impl Delta {
    fn value(&self) -> usize {
        match self {
            Self::N => 0,
            Self::E => 1,
            Self::S => 2,
            Self::W => 3,
        }
    }

    fn turn(&self) -> Self {
        match self {
            Self::N => Self::E,
            Self::E => Self::S,
            Self::S => Self::W,
            Self::W => Self::N,
        }
    }
}

fn get_pos_id(pos: Pos, w: usize) -> usize {
    pos.y * w + pos.x
}

fn get_state_id(pos: Pos, w: usize, delta: Delta) -> usize {
    get_pos_id(pos, w) * 4 + delta.value()
}

fn out_bounds(pos: Pos, bounds: Pos) -> bool {
    pos.x >= bounds.x || pos.y >= bounds.y
}
