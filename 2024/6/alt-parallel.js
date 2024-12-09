const fs = require("fs").promises;
const path = require("path");
const WorkerPool = require("./worker/worker-pool");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = (rows) => {
  const grid = [];
  for (let row of rows) {
    grid.push(row.split(""));
  }
  return grid;
};

const positionsMapping = {
  "^": [-1, 0],
  "<": [0, -1],
  ">": [0, 1],
  v: [1, 0],
};

const moveMapping = {
  "^": ">",
  "<": "^",
  ">": "v",
  v: "<",
};

const part1 = (grid) => {
  let { obstacles, areaLimits, guardPosition, guardDirection } = init(grid);

  const visited = new Set();

  while (!areaLimits.has(`${guardPosition[0]}_${guardPosition[1]}`)) {
    visited.add(`${guardPosition[0]}_${guardPosition[1]}`);

    const nextPosition = [
      guardPosition[0] + positionsMapping[guardDirection][0],
      guardPosition[1] + positionsMapping[guardDirection][1],
    ];

    if (obstacles.has(`${nextPosition[0]}_${nextPosition[1]}`)) {
      guardDirection = moveMapping[guardDirection];
    }

    guardPosition = [
      guardPosition[0] + positionsMapping[guardDirection][0],
      guardPosition[1] + positionsMapping[guardDirection][1],
    ];
  }

  visited.add(`${guardPosition[0]}_${guardPosition[1]}`);

  console.log(visited.size);

  return visited.size;
};

const part2 = async (grid) => {
  let { obstacles, areaLimits, guardPosition, guardDirection } = init(grid);
  let startingPosition = guardPosition;
  let startingDirection = guardDirection;

  const visited = new Set();
  const guardMoves = new Map();

  while (!areaLimits.has(`${guardPosition[0]}_${guardPosition[1]}`)) {
    const key = `${guardPosition[0]}_${guardPosition[1]}`;
    visited.add(key);
    guardMoves.set(key, (guardMoves.get(key) || new Set()).add(guardDirection));

    const nextPosition = [
      guardPosition[0] + positionsMapping[guardDirection][0],
      guardPosition[1] + positionsMapping[guardDirection][1],
    ];

    if (obstacles.has(`${nextPosition[0]}_${nextPosition[1]}`)) {
      guardDirection = moveMapping[guardDirection];
    }

    guardPosition = [
      guardPosition[0] + positionsMapping[guardDirection][0],
      guardPosition[1] + positionsMapping[guardDirection][1],
    ];
  }
  visited.add(`${guardPosition[0]}_${guardPosition[1]}`);
  guardMoves.set(
    `${guardPosition[0]}_${guardPosition[1]}`,
    (
      guardMoves.get(`${guardPosition[0]}_${guardPosition[1]}`) || new Set()
    ).add(guardDirection)
  );

  const pool = WorkerPool(
    Math.max(8, Math.floor(require("os").cpus().length / 2)),
    path.join(__dirname, "./worker/simulation-worker.js")
  );
  const newObstructions = new Set();
  const promises = [];

  for (const cell of visited) {
    const [row, col] = cell.split("_").map(Number);
    if (row === startingPosition[0] && col === startingPosition[1]) continue;

    promises.push(
      new Promise((resolve) => {
        resolve(
          pool.runTask({
            start: startingPosition,
            direction: startingDirection,
            newObstacle: cell,
            areaLimits: Array.from(areaLimits),
            obstacles: Array.from(obstacles),
          })
        );
      })
    );
  }

  const results = await Promise.all(promises);
  pool.terminate();

  results.forEach((result, index) => {
    if (result) {
      const cell = Array.from(visited)[index];
      newObstructions.add(cell);
    }
  });

  console.log(newObstructions.size);

  return newObstructions.size;
};

const init = (grid) => {
  const obstacles = new Set();
  const areaLimits = new Set();
  let guardPosition = [0, 0];
  let guardDirection = "^";

  for (let i = 0; i < grid.length; i++) {
    for (let j = 0; j < grid[i].length; j++) {
      const cell = grid[i][j];
      if (cell === "#") obstacles.add(`${i}_${j}`);
      else if (
        i === grid.length - 1 ||
        j === grid.length - 1 ||
        i === 0 ||
        j === 0
      )
        areaLimits.add(`${i}_${j}`);

      if (["^", "<", ">", "v"].includes(cell)) {
        guardPosition = [i, j];
        guardDirection = cell;
      }
    }
  }
  return { obstacles, areaLimits, guardDirection, guardPosition };
};

readData()
  .then(formatData)
  .then(async (formattedData) => {
    part1(formattedData);
    part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
