use advent2024::bitset::BitSet;
use std::cmp::Ordering;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut count1 = 0;
    let mut count2 = 0;
    let mut orders = [[0 as u8; 100]; 100];
    let mut seen = BitSet::new(100);

    let mut first_half = true;
    for line in reader.lines() {
        let line = line?;
        if first_half {
            if line.len() == 0 {
                first_half = false;
                continue;
            }
            let mut f = line.splitn(2, "|");
            let num1 = f.next().ok_or("Missing first num")?.parse::<u8>()?;
            let num2 = f.next().ok_or("Missing second num")?.parse::<u8>()?;
            let current = orders[num1 as usize][0] + 1;
            orders[num1 as usize][current as usize] = num2;
            orders[num1 as usize][0] = current;
            continue;
        }
        let mut in_order = true;
        let mut nums = line
            .split(",")
            .map(|v| v.parse::<u8>())
            .collect::<Result<Vec<_>, _>>()?;
        seen.reset();
        'outer: for &i in nums.iter() {
            for &o in &orders[i as usize][1..((orders[i as usize][0] + 1) as usize)] {
                if seen.contains(o as usize) {
                    in_order = false;
                    break 'outer;
                }
                seen.insert(i as usize);
            }
        }
        if in_order {
            count1 += nums[nums.len() / 2] as u32;
        } else {
            nums.sort_by(|&a, &b| {
                if a == b {
                    return Ordering::Equal;
                }
                for &i in &orders[a as usize][1..((orders[a as usize][0] + 1) as usize)] {
                    if i == b {
                        return Ordering::Less;
                    }
                }
                for &i in &orders[b as usize][1..((orders[b as usize][0] + 1) as usize)] {
                    if i == a {
                        return Ordering::Greater;
                    }
                }
                Ordering::Equal
            });
            count2 += nums[nums.len() / 2] as u32;
        }
    }

    println!("Part 1: {}", count1);
    println!("Part 2: {}", count2);
    Ok(())
}
