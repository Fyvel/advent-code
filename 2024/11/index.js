const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = (rows) => rows[0].split(" ").map(Number);

const processStone = (stone) => {
  if (stone === 0) return [1];

  const strValue = stone.toString();
  if (strValue.length % 2 === 0) {
    const middle = Math.floor(strValue.length / 2);
    const leftValue = parseInt(strValue.slice(0, middle));
    const rightValue = parseInt(strValue.slice(middle));
    return [leftValue, rightValue];
  }

  return [stone * 2024];
};

const part1 = (stones) => {
  let currentStones = [...stones];
  const depth = 25;

  for (let i = 0; i <= depth; i++) {
    const nextStones = [];
    for (const stone of currentStones) {
      nextStones.push(...processStone(stone));
    }
    currentStones = nextStones;
  }

  console.log(`${currentStones.length}`);
  return currentStones.length;
};

function calculateStoneSize(stone, depth, memo = new Map()) {
  if (depth === 0) return 1;

  const key = `${stone}_${depth}`;
  if (memo.has(key)) return memo.get(key);

  const nextStones = processStone(stone);
  let sum = 0;
  for (const nextStone of nextStones) {
    sum += calculateStoneSize(nextStone, depth - 1, memo);
  }

  memo.set(key, sum);
  return sum;
}

const part2 = (stones) => {
  const depth = 75;
  const memo = new Map();
  let sum = 0;
  for (let i = 0; i < stones.length; i++) {
    sum += calculateStoneSize(stones[i], depth, memo);
  }
  console.log(`${sum}`);
  return sum;
};

readData()
  .then(formatData)
  .then(async (formattedData) => {
    part1(formattedData);
    part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
