use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut count1 = 0;
    let mut count2 = 0;

    let mut row = Vec::new();
    for line in reader.lines() {
        let line = line?;
        row.clear();
        for i in line.split_whitespace() {
            row.push(i.parse::<i8>()?);
        }
        if is_safe(&row, -1) {
            count1 += 1;
            count2 += 1;
        } else {
            for i in 0..(row.len() as i8) {
                if is_safe(&row, i) {
                    count2 += 1;
                    break;
                }
            }
        }
    }

    println!("Part 1: {}", count1);
    println!("Part 2: {}", count2);
    Ok(())
}

fn is_safe(row: &[i8], exclude: i8) -> bool {
    let mut first = true;
    let mut second = true;
    let mut prev = 0;
    let mut inc = true;
    for (n, &num) in row.iter().enumerate() {
        if n as i8 == exclude {
            continue;
        }
        if first {
            first = false;
            prev = num;
            continue;
        }
        let delta = num - prev;
        prev = num;
        let a = delta.abs();
        if a < 1 || a > 3 {
            return false;
        }
        let pos = delta > 0;
        if second {
            second = false;
            inc = pos;
            continue;
        }
        if pos != inc {
            return false;
        }
    }
    true
}
