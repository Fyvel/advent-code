const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = async (rows) => rows;

const part1 = (grid) => {
  let count = 0

  const directions = [
    [-1, -1],
    [-1, 0],
    [-1, 1],
    [0, 1],
    [0, -1],
    [1, -1],
    [1, 0],
    [1, 1]
  ]

  for (let row = 0; row < grid.length; row++)
    for (let col = 0; col < grid[0].length; col++)
      for (let direction of directions)
        if (dfs(grid, row, col, direction, "XMAS"))
          count++

  console.log(count)
  return count
};

const dfs = (grid, row, col, direction, word) => {
  let r = row
  let c = col
  for (let letter of word) {
    // out of bounds
    if (r < 0 || r >= grid.length || c < 0 || c >= grid[0].length)
      return false;
    // wrong letter
    if (grid[r][c] !== letter)
      return false;
    // search next letter
    r += direction[0]
    c += direction[1]
  }
  return true
}

const part2 = (grid) => {
  let count = 0

  const directions = [
    [1, 1], // right down
    [1, -1], // left down
  ]

  for (let row = 0; row < grid.length; row++) {
    for (let col = 0; col < grid[0].length; col++) {
      if (!["M", "S"].includes(grid[row][col])) continue

      const word = grid[row][col] === "M" ? "MAS" : "SAM"
      const middle = Math.floor(word.length / 2)
      const crossingWord = grid[row][col + (word.length - 1)] === "M" ? "MAS" : "SAM"

      if (dfs(grid, row, col, directions[0], word))
        if (dfs(grid, row, col + (word.length - 1), directions[1], crossingWord))
          if (grid[row + middle][col + middle] === word[middle])
            count++
    }
  }

  console.log(count)
  return count
};

readData()
  .then(formatData)
  .then(async (formattedData) => {
    part1(formattedData);
    part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
