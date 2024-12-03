use std::fs;
use std::path::Path;
use std::error::Error;
use tokio;
use regex::Regex;

async fn read_data() -> Result<Vec<String>, Box<dyn Error>> {
    let contents = fs::read_to_string(Path::new("data.txt"))?;
    Ok(contents.lines().map(String::from).collect())
}

async fn format_data(rows: Vec<String>) -> Vec<String> {
		rows
}

fn part1(data: &[String]) -> i32 {
	let mut sum = 0;
    let re = Regex::new(r"mul\((\d{1,3}),(\d{1,3})\)").unwrap();
    
    for line in data {
        for cap in re.captures_iter(line) {
            let a: i32 = cap[1].parse().unwrap();
            let b: i32 = cap[2].parse().unwrap();
            sum += a * b;
        }
    }
    
    println!("sum: {}", sum);
    sum
}

fn part2(data: &[String]) -> i32 {
	let mut sum = 0;
	let re_all = Regex::new(r"(do(n't)?\(\))|(mul\(\d{1,3},\d{1,3}\))").unwrap();
	let re_nums = Regex::new(r"\d+").unwrap();
	
	let mut instruction = String::from("do()");
	
	let multipliers: Vec<String> = data.iter()
			.flat_map(|x| re_all.find_iter(x))
			.map(|m| m.as_str().to_string())
			.collect();
	
	for mul in multipliers {
			if mul.contains("do") {
					instruction = mul.clone();
			} else {
					if instruction != "do()" {
							continue;
					}
					let nums: Vec<i32> = re_nums.find_iter(&mul)
							.map(|m| m.as_str().parse::<i32>().unwrap())
							.collect();
					sum += nums[0] * nums[1];
			}
	}
	
	println!("sum: {}", sum);
	sum
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
	let data = read_data().await?;
	let formatted_data = format_data(data).await;

  part1(&formatted_data);
	part2(&formatted_data);
	
	Ok(())
}

