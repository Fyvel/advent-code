use std::fs::File;
use std::io::{self, BufRead};
use std::path::Path;

fn read_data() -> io::Result<Vec<String>> {
    let path = Path::new("data.txt");
    let file = File::open(&path)?;
    let reader = io::BufReader::new(file);

    let mut data = Vec::new();
    for line in reader.lines() {
        data.push(line?);
    }

    Ok(data)
}

fn format_data(rows: Vec<String>) -> Vec<Vec<char>> {
    rows.into_iter().map(|row| row.chars().collect()).collect()
}

fn part1(grid: &[Vec<char>]) -> usize {
    let mut count = 0;
    let directions = [
        (-1, -1), (-1, 0), (-1, 1),
        (0, 1), (0, -1),
        (1, -1), (1, 0), (1, 1),
    ];

    for row in 0..grid.len() {
        for col in 0..grid[0].len() {
            for direction in &directions {
                if dfs(grid, row as isize, col as isize, *direction, "XMAS") {
                    count += 1;
                }
            }
        }
    }

    println!("{}", count);
    count
}

fn dfs(grid: &[Vec<char>], mut row: isize, mut col: isize, direction: (isize, isize), word: &str) -> bool {
    for letter in word.chars() {
        if row < 0 || row >= grid.len() as isize || col < 0 || col >= grid[0].len() as isize {
            return false;
        }
        if grid[row as usize][col as usize] != letter {
            return false;
        }
        row += direction.0;
        col += direction.1;
    }
    true
}

fn part2(grid: &[Vec<char>]) -> usize {
    let mut count = 0;
    let directions = [
        (1, 1),  // right down
        (1, -1), // left down
    ];

    for row in 0..grid.len() {
        for col in 0..grid[0].len() {
            if grid[row][col] != 'M' && grid[row][col] != 'S' {
                continue;
            }

            let word = if grid[row][col] == 'M' { "MAS" } else { "SAM" };
            let middle = word.len() / 2;

            if col + word.len() - 1 >= grid[0].len() {
                continue;
            }

            let crossing_word = if grid[row][col + word.len() - 1] == 'M' { "MAS" } else { "SAM" };

            if dfs(grid, row as isize, col as isize, directions[0], word) &&
                dfs(grid, row as isize, (col + word.len() - 1) as isize, directions[1], crossing_word) &&
                grid[row + middle][col + middle] == word.chars().nth(middle).unwrap() {
                count += 1;
            }
        }
    }

    println!("{}", count);
    count
}

fn main() -> io::Result<()> {
    let data = read_data()?;
    let formatted_data = format_data(data);
    part1(&formatted_data);
    part2(&formatted_data);
    Ok(())
}