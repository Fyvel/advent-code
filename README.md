# advent-code

[Advent of Code](https://adventofcode.com/)

let's do this

## Run the kikimeter

```sh
# sh compare.sh <day> <number of runs (optional)>
sh compare.sh 1 10
```

## Template

### Node

```js
const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = async (rows) => {};

const part1 = () => {};

const part2 = () => {};

readData()
  .then(formatData)
  .then(async (formattedData) => {
    part1(formattedData);
    part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
```

### Go

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func readData() ([]string, error) {
	data, err := os.ReadFile(filepath.Join(".", "data.txt"))
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

func formatData(rows []string) { }

func part1(){}

func part2(){}

func main() {
	reports, err := readData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	formattedData := formatData(reports)
	part1(formattedData)
	part2(formattedData)
}
```

### Rust

```rs
use std::fs;
use std::path::Path;
use std::error::Error;
use tokio;

async fn read_data() -> Result<Vec<String>, Box<dyn Error>> {
    let contents = fs::read_to_string(Path::new("data.txt"))?;
    Ok(contents.lines().map(String::from).collect())
}

async fn format_data(rows: Vec<String>) -> () {}

fn part1() -> () {}

fn part2() -> () {}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
	let data = read_data().await?;
	let formatted = format_data(data).await;
	part1(formatted.clone());
	part2(formatted);
	Ok(())
}
```
