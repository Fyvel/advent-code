use std::fs;
use std::path::Path;
use std::error::Error;
use std::collections::{HashMap, HashSet};
use tokio;

async fn read_data() -> Result<Vec<String>, Box<dyn Error>> {
	let contents = fs::read_to_string(Path::new("data.txt"))?;
	Ok(contents.lines().map(String::from).collect())
}


#[derive(Debug)]
struct FormattedData {
    rules: Vec<String>,
    updates: Vec<String>,
}

fn format_data(rows: Vec<String>) -> FormattedData {
    let mut rules = Vec::new();
    let mut updates = Vec::new();
    let mut is_end_of_rules = false;

    for line in rows.iter() {
        if line.is_empty() {
            is_end_of_rules = true;
            continue;
        }
        if !is_end_of_rules {
            rules.push(line.to_string());
        } else {
            updates.push(line.to_string());
        }
    }

    FormattedData { rules, updates }
}

fn part1(data: &FormattedData) -> i32 {
    let rules_set: HashSet<_> = data.rules.iter().cloned().collect();
    let mut valid_updates = Vec::new();

    for update in &data.updates {
        let pages: Vec<&str> = update.split(',').collect();
        let mut is_valid_update = true;

        for i in 0..pages.len() - 1 {
            let failing_rule = format!("{}|{}", pages[i + 1], pages[i]);
            if rules_set.contains(&failing_rule) {
                is_valid_update = false;
                break;
            }
        }
        if is_valid_update {
            valid_updates.push(pages);
        }
    }

    let mut sum = 0;
    for valid_update in valid_updates {
        sum += valid_update[valid_update.len() / 2].parse::<i32>().unwrap();
    }

    println!("{}", sum);
    sum
}

fn part2(data: &FormattedData) -> i32 {
    let mut corrected_updates = Vec::new();

    for update in &data.updates {
        let pages: Vec<i32> = update.split(',').map(|x| x.parse().unwrap()).collect();

        let mut rules_set = Vec::new();
        for rule in &data.rules {
            let parts: Vec<i32> = rule.split('|').map(|x| x.parse().unwrap()).collect();
            if pages.contains(&parts[0]) && pages.contains(&parts[1]) {
                rules_set.push((parts[0], parts[1]));
            }
        }

        let mut dependencies: HashMap<i32, i32> = HashMap::new();
        for &(_, b) in &rules_set {
            *dependencies.entry(b).or_insert(0) += 1;
        }

        let mut corrected = Vec::new();
        while corrected.len() < pages.len() {
            for &page in &pages {
                if !corrected.contains(&page) {
                    if dependencies.get(&page).unwrap_or(&0) < &1 {
                        corrected.push(page);
                        for &(a, b) in &rules_set {
                            if page == a {
                                if let Some(count) = dependencies.get_mut(&b) {
                                    *count -= 1;
                                }
                            }
                        }
                    }
                }
            }
        }

        let original = pages
            .iter()
            .map(|x| x.to_string())
            .collect::<Vec<String>>()
            .join(",");

        let corrected_str = corrected
            .iter()
            .map(|x| x.to_string())
            .collect::<Vec<String>>()
            .join(",");

        if original != corrected_str {
            corrected_updates.push(corrected);
        }
    }

    let mut sum = 0;
    for corrected_update in corrected_updates {
        sum += corrected_update[corrected_update.len() / 2];
    }

    println!("{}", sum);
    sum
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
	let data = read_data().await?;
	let formatted = format_data(data);
	part1(&formatted);
	part2(&formatted);
	Ok(())
}