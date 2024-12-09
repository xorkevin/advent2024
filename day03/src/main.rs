use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mul_instr_regex = Regex::new(r"mul\((\d{1,3}),(\d{1,3})\)|do(?:n't)?\(\)").unwrap();

    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut sum1 = 0;
    let mut sum2 = 0;
    let mut enabled = true;

    for line in reader.lines() {
        let line = line?;
        for i in mul_instr_regex.captures_iter(&line) {
            if i[0].starts_with("don") {
                enabled = false;
                continue;
            }
            if i[0].starts_with("d") {
                enabled = true;
                continue;
            }
            let product = i[1].parse::<u32>()? * i[2].parse::<u32>()?;
            sum1 += product;
            if enabled {
                sum2 += product;
            }
        }
    }

    println!("Part 1: {}", sum1);
    println!("Part 2: {}", sum2);
    Ok(())
}
