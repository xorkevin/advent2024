use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut count1 = 0;
    let mut count2 = 0;

    for line in reader.lines() {
        let line = line?;
        let mut ff = line.splitn(2, ": ");
        let target = ff.next().ok_or("Invalid line")?.parse::<u64>()?;
        let nums = ff
            .next()
            .ok_or("Invalid line")?
            .split_whitespace()
            .map(|v| v.parse::<u64>())
            .collect::<Result<Vec<_>, _>>()?;
        if search(target, nums[0], &nums[1..], false) {
            count1 += target;
        } else if search(target, nums[0], &nums[1..], true) {
            count2 += target;
        }
    }

    println!("Part 1: {}", count1);
    println!("Part 2: {}", count1 + count2);
    Ok(())
}

fn search(target: u64, start: u64, rest: &[u64], pt2: bool) -> bool {
    if rest.len() == 0 {
        return target == start;
    }
    if search(target, start + rest[0], &rest[1..], pt2) {
        return true;
    }
    if search(target, start * rest[0], &rest[1..], pt2) {
        return true;
    }
    if !pt2 {
        return false;
    }
    search(
        target,
        start * magnitude(rest[0]) + rest[0],
        &rest[1..],
        pt2,
    )
}

fn magnitude(mut i: u64) -> u64 {
    let mut mag = 10;
    while i >= 10 {
        i /= 10;
        mag *= 10;
    }
    return mag;
}
