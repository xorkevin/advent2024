use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut nums1 = Vec::new();
    let mut nums2 = Vec::new();
    for line in reader.lines() {
        let line = line?;
        let mut fields = line.split_whitespace();
        let num1 = fields.next().ok_or("Invalid")?.parse::<usize>()?;
        let num2 = fields.next().ok_or("Invalid")?.parse::<usize>()?;
        nums1.push(num1);
        nums2.push(num2);
    }
    nums1.sort();
    nums2.sort();

    let min_num2 = nums2[0];
    let max_num2 = nums2[nums2.len() - 1];

    let mut counts = vec![0; max_num2 - min_num2 + 1];
    for i in nums2.iter() {
        counts[i - min_num2] += 1;
    }

    let (s1, s2) = nums1
        .iter()
        .zip(nums2.iter())
        .fold((0, 0), |(s1, s2), (&v1, &v2)| {
            (
                s1 + v1.abs_diff(v2),
                s2 + if v1 < min_num2 {
                    0
                } else {
                    v1 * counts[v1 - min_num2]
                },
            )
        });

    println!("Part 1: {}", s1);
    println!("Part 2: {}", s2);
    Ok(())
}
