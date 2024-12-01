const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = async (rows) => {
  const left = [];
  const right = [];
  rows.forEach((row) => {
    const [a, b] = row.split("   ");
    left.push(a);
    right.push(b);
  });
  return [left, right];
};

const part1 = async ([left, right]) => {
  left.sort();
  right.sort();
  let distance = 0;
  for (let i = 0; i < left.length; i++) {
    distance += Math.abs(left[i] - right[i]);
  }
  console.log({ distance });
  return distance;
};

const part2 = async ([left, right]) => {
  let similarity = 0;
  const map = new Map();
  for (let i = 0; i < right.length; i++) {
    map.set(right[i], (map.get(right[i]) || 0) + 1);
  }
  for (let i = 0; i < left.length; i++) {
    similarity += left[i] * (map.get(left[i]) || 0);
  }
  console.log({ similarity });
  return similarity;
};

readData()
  .then(formatData)
  .then(async (formattedData) => {
    await part1(formattedData);
    await part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
