const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = async (rows) => rows.map((row) => row.split("").map(Number));

function bfs(grid, startPoints, aggregator, visited, scores) {
  let queue = [...startPoints];
  const directions = [
    [-1, 0],
    [1, 0],
    [0, 1],
    [0, -1],
  ];

  while (queue.length) {
    const cell = queue.shift();
    const x = cell[0];
    const y = cell[1];
    const cellKey = `${x}_${y}`;
    visited[cellKey] = true;

    aggregator(x, y, grid[x][y], scores);

    for (const dir of directions) {
      const { newX, newY } = { newX: x + dir[0], newY: y + dir[1] };

      if (newX < 0 || newY < 0 || newX >= grid.length || newY >= grid[y].length)
        continue;

      const neighbour = grid[newX][newY];
      const key = `${newX}_${newY}`;

      if (neighbour == grid[x][y] + 1) {
        queue.push([newX, newY]);
        visited[key] = true;
      }
    }
  }
}

const part1 = (grid) => {
  const startPoints = [];
  for (let i = 0; i < grid.length; i++)
    for (let j = 0; j < grid[i].length; j++)
      if (grid[i][j] === 0) startPoints.push([i, j]);

  const scores = new Map();
  for (const startPoint of startPoints) {
    const cellVisited = {};
    const summitSeen = {};

    const aggregator = (x, y, value) => {
      const key = `${x}_${y}`;
      if (value === 9 && !summitSeen[key]) {
        trailHeadKey = `${startPoint[0]}_${startPoint[1]}`;
        summitSeen[key] = true;
        scores.set(trailHeadKey, (scores.get(trailHeadKey) || 0) + 1);
      }
    };

    bfs(grid, [startPoint], aggregator, cellVisited, scores);
  }

  let sum = 0;
  for (const score of scores.values()) {
    sum += score;
  }
  console.log(sum);
  return sum;
};

const part2 = (grid) => {
  const startPoints = [];
  for (let i = 0; i < grid.length; i++)
    for (let j = 0; j < grid[i].length; j++)
      if (grid[i][j] === 0) startPoints.push([i, j]);

  const scores = new Map();
  for (const startPoint of startPoints) {
    const cellVisited = {};

    const aggregator = (x, y, value) => {
      if (value === 9) {
        trailHeadKey = `${startPoint[0]}_${startPoint[1]}`;
        scores.set(trailHeadKey, (scores.get(trailHeadKey) || 0) + 1);
      }
    };

    bfs(grid, [startPoint], aggregator, cellVisited, scores);
  }

  let sum = 0;
  for (const score of scores.values()) {
    sum += score;
  }
  console.log(sum);
  return sum;
};

readData()
  .then(formatData)
  .then(async (formattedData) => {
    part1(formattedData);
    part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
