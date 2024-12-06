const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = (rows) => {
  const grid = []
  for (let row of rows) {
    grid.push(row.split(""));
  }
  return grid;
};

const part1 = (grid) => {
  const obstacles = new Set()
  const areaLimits = new Set()
  let guardPosition = [0, 0];
  let guardDirection = "^";
  const positionsMapping = {
    "^": [-1, 0],
    "<": [0, -1],
    ">": [0, 1],
    "v": [1, 0],
  }
  const moveMapping = {
    "^": ">",
    "<": "^",
    ">": "v",
    "v": "<",
  }

  for (let i = 0; i < grid.length; i++) {
    for (let j = 0; j < grid[i].length; j++) {
      const cell = grid[i][j];
      if (cell === "#")
        obstacles.add(`${i}_${j}`)
      else if (i === grid.length - 1 || j === grid.length - 1 || i === 0 || j === 0)
        areaLimits.add(`${i}_${j}`)

      if (["^", "<", ">", "v"].includes(cell)) {
        guardPosition = [i, j];
        guardDirection = cell;
      }
    }
  }

  const visited = new Set();

  while (!areaLimits.has(`${guardPosition[0]}_${guardPosition[1]}`)) {
    visited.add(`${guardPosition[0]}_${guardPosition[1]}`);

    const nextPosition = [
      guardPosition[0] + positionsMapping[guardDirection][0],
      guardPosition[1] + positionsMapping[guardDirection][1]
    ]

    if (obstacles.has(`${nextPosition[0]}_${nextPosition[1]}`)) {
      guardDirection = moveMapping[guardDirection]
    }

    guardPosition = [
      guardPosition[0] + positionsMapping[guardDirection][0],
      guardPosition[1] + positionsMapping[guardDirection][1]
    ]
  }

  visited.add(`${guardPosition[0]}_${guardPosition[1]}`);

  console.log(visited.size)

  return visited.size
};

const part2 = (grid) => {
  const obstacles = new Set()
  const areaLimits = new Set()
  let guardPosition = [0, 0];
  let guardDirection = "^";
  const positionsMapping = {
    "^": [-1, 0],
    "<": [0, -1],
    ">": [0, 1],
    "v": [1, 0],
  }
  const moveMapping = {
    "^": ">",
    "<": "^",
    ">": "v",
    "v": "<",
  }

  for (let i = 0; i < grid.length; i++) {
    for (let j = 0; j < grid[i].length; j++) {
      const cell = grid[i][j];
      if (cell === "#")
        obstacles.add(`${i}_${j}`)

      if (cell !== "#" && (i === grid.length - 1 || j === grid[0].length - 1 || i === 0 || j === 0))
        areaLimits.add(`${i}_${j}`)

      if (["^", "<", ">", "v"].includes(cell)) {
        guardPosition = [i, j];
        guardDirection = cell;
      }
    }
  }

  const visited = new Set();
  const newObstructions = new Set()

  const simulate = (startRow, startCol, startDirection, areaLimits, obstacles, visited) => {
    let row = startRow;
    let col = startCol;
    let direction = startDirection;
    const localVisited = new Set()

    while (true) {
      const positionKey = `${row}_${col}_${direction}`;

      if (visited.has(positionKey) || localVisited.has(positionKey)) {
        return true;
      }
      localVisited.add(positionKey);

      if (areaLimits.has(`${row}_${col}`)) {
        return false;
      }

      if (obstacles.has(`${row}_${col}`)) {
        direction = moveMapping[direction];
        row += positionsMapping[direction][0];
        col += positionsMapping[direction][1];
      } else {
        row += positionsMapping[direction][0];
        col += positionsMapping[direction][1];
      }
    }
  };

  while (!areaLimits.has(`${guardPosition[0]}_${guardPosition[1]}`)) {
    visited.add(`${guardPosition[0]}_${guardPosition[1]}_${guardDirection}`)

    const nextPosition = [
      guardPosition[0] + positionsMapping[guardDirection][0],
      guardPosition[1] + positionsMapping[guardDirection][1]
    ]

    if (obstacles.has(`${nextPosition[0]}_${nextPosition[1]}`)) {
      guardDirection = moveMapping[guardDirection]
    }

    if (simulate(guardPosition[0], guardPosition[1], moveMapping[guardDirection], areaLimits, obstacles, visited)) {
      newObstructions.add(`${nextPosition[0]}_${nextPosition[1]}`)
    }
    guardPosition = [
      guardPosition[0] + positionsMapping[guardDirection][0],
      guardPosition[1] + positionsMapping[guardDirection][1]
    ]
  }

  console.log(newObstructions.size)

  return newObstructions.size
};



readData()
  .then(formatData)
  .then(async (formattedData) => {
    part1(formattedData);
    part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
