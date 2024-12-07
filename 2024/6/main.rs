use std::fs;
use std::path::Path;
use std::error::Error;
use std::collections::{HashMap, HashSet};
use tokio;

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

fn part1(grid: Grid) -> usize {
    let mut state = init(&grid);
    let mut visited = HashSet::new();

    while !state.area_limits.contains(&format!("{}_{}", state.guard_position.0, state.guard_position.1)) {
        visited.insert(format!("{}_{}", state.guard_position.0, state.guard_position.1));

        let offset = get_position_offset(state.guard_direction);
        let next_position = (
            state.guard_position.0 + offset.0,
            state.guard_position.1 + offset.1,
        );

        if state.obstacles.contains(&format!("{}_{}", next_position.0, next_position.1)) {
            state.guard_direction = get_next_direction(state.guard_direction);
        }

        let offset = get_position_offset(state.guard_direction);
        state.guard_position = (
            state.guard_position.0 + offset.0,
            state.guard_position.1 + offset.1,
        );
    }

    visited.insert(format!("{}_{}", state.guard_position.0, state.guard_position.1));
    println!("Part 1: {}", visited.len());
    visited.len()
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

fn part2(grid: Grid) -> usize {
    let state = init(&grid);
    let start_pos = state.guard_position;
    let start_dir = state.guard_direction;
    let mut visited = HashSet::new();
    let mut guard_moves: HashMap<String, HashSet<char>> = HashMap::new();
    let mut current_state = state.clone();

    while !current_state.area_limits.contains(&format!("{}_{}", current_state.guard_position.0, current_state.guard_position.1)) {
        let key = format!("{}_{}", current_state.guard_position.0, current_state.guard_position.1);
        visited.insert(key.clone());
        guard_moves
            .entry(key)
            .or_insert_with(HashSet::new)
            .insert(current_state.guard_direction);

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

    let mut new_obstructions = HashSet::new();
    for cell in visited {
        let coords: Vec<i32> = cell
            .split('_')
            .map(|x| x.parse().unwrap())
            .collect();
        
        if (coords[0], coords[1]) == start_pos {
            continue;
        }

        if simulate(
            start_pos,
            start_dir,
            &cell,
            &state.area_limits,
            &state.obstacles,
        ) {
            new_obstructions.insert(cell);
        }
    }

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