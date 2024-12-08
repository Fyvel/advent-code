const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = (rows) => rows.map((row) => row.split(""));

// helper function to visualise the grid
const visualiseGrid = (grid, antiNodes = new Set()) => {
  console.log(
    grid
      .map((row, i) =>
        row
          .map((cell, j) => {
            if (antiNodes.has(`${j}_${i}`)) {
              return `◆`;
            }
            return cell;
          })
          .join("")
      )
      .join(" \n")
  );
  console.log(Array(grid[0].length).fill("-").join(""));
};

const part1 = (grid) => {
  visualiseGrid(grid);

  const antennasMap = {};
  const antiNodes = new Set();

  for (let i = 0; i < grid.length; i++) {
    for (let j = 0; j < grid[i].length; j++) {
      const cell = grid[i][j];

      if (!cell.match(/[a-zA-Z0-9]/)) continue;

      if (!antennasMap[cell]) antennasMap[cell] = [];
      antennasMap[cell].push({ x: j, y: i });

      for (const frequency in antennasMap) {
        const locations = antennasMap[frequency];

        for (let i = 0; i < locations.length; i++) {
          for (let j = i + 1; j < locations.length; j++) {
            const location1 = locations[i];
            const location2 = locations[j];

            const dx = location1.x - location2.x;
            const dy = location1.y - location2.y;

            const antiNode1 = [location1.x + dx, location1.y + dy];
            const antiNode2 = [location2.x - dx, location2.y - dy];

            if (
              antiNode1[0] >= 0 &&
              antiNode1[1] >= 0 &&
              antiNode1[0] < grid[0].length &&
              antiNode1[1] < grid.length
            ) {
              antiNodes.add(`${antiNode1[0]}_${antiNode1[1]}`);
            }

            if (
              antiNode2[0] >= 0 &&
              antiNode2[1] >= 0 &&
              antiNode2[0] < grid[0].length &&
              antiNode2[1] < grid.length
            ) {
              antiNodes.add(`${antiNode2[0]}_${antiNode2[1]}`);
            }
          }
        }
      }
    }
  }

  visualiseGrid(grid, antiNodes);

  console.log(antiNodes.size);
  return antiNodes.size;
};

const part2 = (grid) => {
  visualiseGrid(grid);

  const antennasMap = {};
  const antiNodes = new Set();

  for (let i = 0; i < grid.length; i++) {
    for (let j = 0; j < grid[i].length; j++) {
      const cell = grid[i][j];

      if (!cell.match(/[a-zA-Z0-9]/)) continue;

      if (!antennasMap[cell]) antennasMap[cell] = [];
      antennasMap[cell].push({ x: j, y: i });

      for (const frequency in antennasMap) {
        const locations = antennasMap[frequency];

        for (let i = 0; i < locations.length; i++) {
          for (let j = i + 1; j < locations.length; j++) {
            const location1 = locations[i];
            const location2 = locations[j];

            const dx = location1.x - location2.x;
            const dy = location1.y - location2.y;

            let withinBounds = true;
            let k = 0;

            while (withinBounds) {
              const localAntiNodes = new Set();
              const antiNode1 = [location1.x + dx * k, location1.y + dy * k];
              const antiNode2 = [location2.x - dx * k, location2.y - dy * k];

              if (
                antiNode1[0] >= 0 &&
                antiNode1[1] >= 0 &&
                antiNode1[0] < grid[0].length &&
                antiNode1[1] < grid.length
              ) {
                localAntiNodes.add(`${antiNode1[0]}_${antiNode1[1]}`);
              }

              if (
                antiNode2[0] >= 0 &&
                antiNode2[1] >= 0 &&
                antiNode2[0] < grid[0].length &&
                antiNode2[1] < grid.length
              ) {
                localAntiNodes.add(`${antiNode2[0]}_${antiNode2[1]}`);
              }

              if (localAntiNodes.size > 0) {
                localAntiNodes.forEach((node) => antiNodes.add(node));
              } else {
                withinBounds = false;
              }
              k++;
            }
          }
        }
      }
    }
  }
  visualiseGrid(grid, antiNodes);

  console.log(antiNodes.size);
  return antiNodes.size;
};

readData()
  .then(formatData)
  .then(async (formattedData) => {
    part1(formattedData);
    part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
