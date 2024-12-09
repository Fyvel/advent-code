use std::fs;
use std::path::Path;
use std::error::Error;
use std::collections::{HashMap, HashSet};
use tokio;
use rayon::prelude::*;

type Grid = Vec<Vec<char>>;
type Position = (i32, i32);

#[derive(Clone)]
struct State {
    obstacles: HashSet<String>,
    area_limits: HashSet<String>,
    guard_position: Position,
    guard_direction: char,
}

async fn read_data() -> Result<Vec<String>, Box<dyn Error>> {
    let contents = fs::read_to_string(Path::new("data.txt"))?;
    Ok(contents.lines().map(String::from).collect())
}

async fn format_data(rows: Vec<String>) -> Grid {
    rows.iter()
        .map(|row| row.chars().collect())
        .collect()
}

fn get_position_offset(direction: char) -> Position {
    match direction {
        '^' => (-1, 0),
        '<' => (0, -1),
        '>' => (0, 1),
        'v' => (1, 0),
        _ => panic!("Invalid direction"),
    }
}

fn get_next_direction(direction: char) -> char {
    match direction {
        '^' => '>',
        '<' => '^',
        '>' => 'v',
        'v' => '<',
        _ => panic!("Invalid direction"),
    }
}

fn init(grid: &Grid) -> State {
    let mut obstacles = HashSet::new();
    let mut area_limits = HashSet::new();
    let mut guard_position = (0, 0);
    let mut guard_direction = '^';

    for (i, row) in grid.iter().enumerate() {
        for (j, &cell) in row.iter().enumerate() {
            let pos = format!("{}_{}", i, j);
            
            if cell == '#' {
                obstacles.insert(pos);
            } else if i == 0 || j == 0 || i == grid.len() - 1 || j == row.len() - 1 {
                area_limits.insert(pos);
            }

            if "^<>v".contains(cell) {
                guard_position = (i as i32, j as i32);
                guard_direction = cell;
            }
        }
    }

    State {
        obstacles,
        area_limits,
        guard_position,
        guard_direction,
    }
}

fn simulate(
    start: Position,
    direction: char,
    new_obstacle: &str,
    area_limits: &HashSet<String>,
    obstacles: &HashSet<String>,
) -> bool {
    let mut local_moves: HashMap<String, HashSet<char>> = HashMap::new();
    let mut pos = start;
    let mut dir = direction;

    loop {
        let key = format!("{}_{}", pos.0, pos.1);
        if let Some(moves) = local_moves.get(&key) {
            if moves.contains(&dir) {
                return true;
            }
        }

        local_moves
            .entry(key)
            .or_insert_with(HashSet::new)
            .insert(dir);

        let offset = get_position_offset(dir);
        let next_pos = (pos.0 + offset.0, pos.1 + offset.1);
        let next_key = format!("{}_{}", next_pos.0, next_pos.1);

        if area_limits.contains(&next_key) && next_key != new_obstacle {
            return false;
        }

        if obstacles.contains(&next_key) || next_key == new_obstacle {
            dir = get_next_direction(dir);
        } else {
            pos = next_pos;
        }
    }
}

fn get_initial_path(state: &State) -> HashSet<String> {
    let mut visited = HashSet::new();
    let mut current_state = state.clone();

    while !state.area_limits.contains(&format!("{}_{}", current_state.guard_position.0, current_state.guard_position.1)) {
        visited.insert(format!("{}_{}", current_state.guard_position.0, current_state.guard_position.1));

        let offset = get_position_offset(current_state.guard_direction);
        let next_position = (
            current_state.guard_position.0 + offset.0,
            current_state.guard_position.1 + offset.1,
        );

        if current_state.obstacles.contains(&format!("{}_{}", next_position.0, next_position.1)) {
            current_state.guard_direction = get_next_direction(current_state.guard_direction);
        }

        let offset = get_position_offset(current_state.guard_direction);
        current_state.guard_position = (
            current_state.guard_position.0 + offset.0,
            current_state.guard_position.1 + offset.1,
        );
    }

    visited.insert(format!("{}_{}", current_state.guard_position.0, current_state.guard_position.1));
    visited
}

fn part1(grid: Grid) -> usize {
    let state = init(&grid);
    let visited = get_initial_path(&state);
    println!("Part 1: {}", visited.len());
    visited.len()
}

fn part2(grid: Grid) -> usize {
    let state = init(&grid);
    let visited = get_initial_path(&state);
    
    // Convert visited to Vec for parallel iteration
    let visited_vec: Vec<String> = visited
        .into_iter()
        .filter(|cell| {
            let coords: Vec<i32> = cell
                .split('_')
                .map(|x| x.parse().unwrap())
                .collect();
            (coords[0], coords[1]) != state.guard_position
        })
        .collect();

    // Parallel simulation
    let new_obstructions: HashSet<String> = visited_vec
        .par_iter()
        .filter(|&cell| {
            simulate(
                state.guard_position,
                state.guard_direction,
                cell,
                &state.area_limits,
                &state.obstacles,
            )
        })
        .cloned()
        .collect();

    println!("Part 2: {}", new_obstructions.len());
    new_obstructions.len()
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let data = read_data().await?;
    let formatted = format_data(data).await;
    part1(formatted.clone());
    part2(formatted);
    Ok(())
}