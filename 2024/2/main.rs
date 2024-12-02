use std::collections::HashMap;
use std::fs::read_to_string;

fn read_data() -> Result<Vec<String>, std::io::Error> {
    let data = read_to_string("data.txt")?;
    Ok(data.lines().map(String::from).collect())
}

fn format_data(reports: &[String]) -> Vec<Vec<i32>> {
    reports
        .iter()
        .map(|report| {
            report
                .split_whitespace()
                .filter_map(|num| num.parse::<i32>().ok())
                .collect()
        })
        .collect()
}

fn part1(levels: &[Vec<i32>]) -> i32 {
    let safe_level_seen = levels.iter().filter(|level| check_level(level)).count();
    println!("safeLevelSeen: {}", safe_level_seen);
    safe_level_seen as i32
}

fn part2(levels: &[Vec<i32>]) -> i32 {
    let mut memo: HashMap<String, bool> = HashMap::new();
    let mut safe_level_seen = 0;

    for level in levels {
        let level_key = format!("{:?}", level);
        if let Some(&is_safe) = memo.get(&level_key) {
            if is_safe {
                safe_level_seen += 1;
            }
            continue;
        }

        let mut is_safe_level = check_level(level);
        memo.insert(level_key.clone(), is_safe_level);

        if !is_safe_level {
            is_safe_level = retry(level, &mut memo);
        }

        if is_safe_level {
            safe_level_seen += 1;
        }
    }
    println!("safeLevelSeen: {}", safe_level_seen);
    safe_level_seen
}

fn check_level(level: &Vec<i32>) -> bool {
    let mut is_safe_level = true;
    let mut direction = String::new();

    for j in 0..level.len() - 1 {
        let current_direction = if level[j] >= level[j + 1] {
            "dec".to_string()
        } else {
            "inc".to_string()
        };

        if !direction.is_empty() && direction != current_direction {
            is_safe_level = false;
            break;
        }
        direction = current_direction;

        let gap = (level[j] - level[j + 1]).abs();
        let in_range = gap >= 1 && gap <= 3;
        if !in_range {
            is_safe_level = false;
            break;
        }
    }
    is_safe_level
}

fn retry(level: &Vec<i32>, memo: &mut HashMap<String, bool>) -> bool {
    let mut is_safe_level = false;

    for j in 0..level.len() {
        let mut new_level: Vec<i32> = level[..j].to_vec();
        new_level.extend(&level[j + 1..]);
        let new_level_key = format!("{:?}", new_level);

        if let Some(&safe) = memo.get(&new_level_key) {
            is_safe_level = safe;
            if is_safe_level {
                break;
            }
        }

        is_safe_level = check_level(&new_level);
        memo.insert(new_level_key, is_safe_level);

        if is_safe_level {
            break;
        }
    }
    is_safe_level
}

fn main() {
    match read_data() {
        Ok(reports) => {
            let formatted_data = format_data(&reports);
            part1(&formatted_data);
            part2(&formatted_data);
        }
        Err(e) => println!("Error: {}", e),
    }
}
