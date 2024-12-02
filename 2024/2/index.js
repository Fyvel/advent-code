const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = async (reports) => {
  const levels = [];
  for (let i = 0; i < reports.length; i++) {
    levels.push(reports[i].split(" ").map(Number));
  }
  return levels;
};

const part1 = (levels) => {
  let safeLevelSeen = 0;
  for (let i = 0; i < levels.length; i++) {
    const level = levels[i];
    let isSafeLevel = true;
    let direction = "";

    for (let j = 0; j < level.length - 1; j++) {
      const currentDirection = level[j] >= level[j + 1] ? "dec" : "inc";
      if (direction && direction !== currentDirection) {
        isSafeLevel = false;
        break;
      }
      direction = currentDirection;

      const gap = Math.abs(level[j] - level[j + 1]);
      const inRange = gap >= 1 && gap <= 3;
      if (!inRange) {
        isSafeLevel = false;
        break;
      }
    }
    if (isSafeLevel) safeLevelSeen++;
  }
  console.log({ safeLevelSeen });
  return safeLevelSeen;
};

const part2 = (levels) => {
  const memo = new Map();
  let safeLevelSeen = 0;

  for (let i = 0; i < levels.length; i++) {
    const level = levels[i];
    const levelKey = level.join("");

    if (memo.has(levelKey)) {
      if (memo.get(levelKey)) safeLevelSeen++;
      continue;
    }

    let isSafeLevel = checkLevel(level);
    memo.set(levelKey, isSafeLevel);

    if (!isSafeLevel) isSafeLevel = retry(level, memo);
    if (isSafeLevel) safeLevelSeen++;
  }
  console.log({ safeLevelSeen });
  return safeLevelSeen;
};

const checkLevel = (level) => {
  let isSafeLevel = true;
  let direction = "";

  for (let j = 0; j < level.length - 1; j++) {
    const currentDirection = level[j] >= level[j + 1] ? "dec" : "inc";
    if (direction && direction !== currentDirection) {
      isSafeLevel = false;
      break;
    }
    direction = currentDirection;

    const gap = Math.abs(level[j] - level[j + 1]);
    const inRange = gap >= 1 && gap <= 3;
    if (!inRange) {
      isSafeLevel = false;
      break;
    }
  }
  return isSafeLevel;
};

const retry = (level, memo) => {
  let isSafeLevel;
  for (let j = 0; j < level.length; j++) {
    const newLevel = [...level.slice(0, j), ...level.slice(j + 1)];
    const newLevelKey = newLevel.join("");

    if (memo.has(newLevelKey)) {
      isSafeLevel = memo.get(newLevelKey);
      if (isSafeLevel) break;
    }

    isSafeLevel = checkLevel(newLevel);
    memo.set(newLevelKey, isSafeLevel);

    if (isSafeLevel) break;
  }
  return isSafeLevel;
};

readData()
  .then(formatData)
  .then((formattedData) => {
    part1(formattedData);
    part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
