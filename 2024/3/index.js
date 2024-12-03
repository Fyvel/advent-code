const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = async (rows) => rows;

const part1 = (data) => {
  let sum = 0;
  const multipliers = data
    .map((x) => x.match(/mul\(\d{1,3},\d{1,3}\)/g))
    .flat();
  for (const mul of multipliers) {
    let [a, b] = mul.match(/(\d)+/g);
    sum += a * b;
  }
  console.log({ sum });
  return sum;
};

const part2 = (data) => {
  let sum = 0;
  const multipliers = data
    .map((x) => x.match(/(do(n't)?\(\))|(mul\(\d{1,3},\d{1,3}\))/g))
    .flat();

  let instruction = "do()";
  for (const mul of multipliers) {
    if (mul.match(/(do(n't)?\(\))/g)) instruction = mul;
    else {
      if (instruction !== "do()") continue;
      let [a, b] = mul.match(/(\d)+/g);
      sum += a * b;
    }
  }
  console.log({ sum });
  return sum;
};

readData()
  .then(formatData)
  .then(async (formattedData) => {
    part1(formattedData);
    part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
