const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

// JS moment...
// const formatData_using_map_which_fucking_doesnt_work = async (rows) => {
//   const equations = new Map();
//   for (let row of rows) {
//     const [testValue, numbers] = row.split(": ");
//     equations.set(Number(testValue), numbers.split(" ").map(Number));
//   }
//   return equations;
// };

const formatData = async (rows) => {
  const equations = rows.map((row) =>
    row.split(": ").map((entries) => entries.split(" ").map(Number))
  );
  return equations;
};

const part1 = (equations) => {
  let sum = 0;

  for (const [target, numbers] of equations) {
    const [first, ...rest] = numbers;
    let evaluations = new Set([first]);

    for (let i = 0; i < rest.length; i++) {
      const localEvaluations = new Set();

      for (const left of evaluations) {
        localEvaluations.add(left * rest[i]);
        localEvaluations.add(left + rest[i]);
      }

      evaluations = localEvaluations;
    }

    if (evaluations.has(Number(target))) {
      sum += Number(target);
    }
  }

  console.log(sum);
  return sum;
};

const part2 = (equations) => {
  let sum = 0;

  for (const [target, numbers] of equations) {
    const [first, ...rest] = numbers;
    let evaluations = new Set([first]);

    for (let i = 0; i < rest.length; i++) {
      const localEvaluations = new Set();

      for (const left of evaluations) {
        localEvaluations.add(left * rest[i]);
        localEvaluations.add(left + rest[i]);
        localEvaluations.add(Number(`${left}${rest[i]}`));
      }

      evaluations = localEvaluations;
    }

    if (evaluations.has(Number(target))) {
      sum += Number(target);
    }
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
