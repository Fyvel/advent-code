use std::fs;
use std::path::Path;
use std::error::Error;
use std::collections::HashMap;
use tokio;

async fn read_data() -> Result<Vec<String>, Box<dyn Error>> {
    let contents = fs::read_to_string(Path::new("data.txt"))?;
    Ok(contents.lines().map(String::from).collect())
}

async fn format_data(rows: Vec<String>) -> (Vec<String>, Vec<String>) {
    let mut left = Vec::new();
    let mut right = Vec::new();

    for row in rows {
        let parts: Vec<&str> = row.split("   ").collect();
        if parts.len() == 2 {
            left.push(parts[0].to_string());
            right.push(parts[1].to_string());
        }
    }

    (left, right)
}

fn part1(data: (Vec<String>, Vec<String>)) -> i32 {
    let (mut left, mut right) = data;
    left.sort();
    right.sort();
    
    let mut distance = 0;
    for i in 0..left.len() {
        let left_num: i32 = left[i].parse().unwrap_or(0);
        let right_num: i32 = right[i].parse().unwrap_or(0);
        distance += (left_num - right_num).abs();
    }
    
    println!("distance: {}", distance);
    distance
}


fn part2(data: (Vec<String>, Vec<String>)) -> i32 {
	let (left, right) = data;
	let mut similarity = 0;
	let mut map = HashMap::new();
	
	for num_str in right {
			let num: i32 = num_str.parse().unwrap_or(0);
			*map.entry(num).or_insert(0) += 1;
	}
	
	for num_str in left {
			let num: i32 = num_str.parse().unwrap_or(0);
			similarity += num * map.get(&num).unwrap_or(&0);
	}
	
	println!("similarity: {}", similarity);
	similarity
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
	let data = read_data().await?;
	let formatted = format_data(data).await;
	part1(formatted.clone());
	part2(formatted);
	Ok(())
}