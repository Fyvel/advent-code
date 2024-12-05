const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = async (rows) => {
  const rules = [];
  const updates = [];
  let isEndOfRules = false;

  for (const row of rows) {
    if (row === "") {
      isEndOfRules = true;
      continue;
    }
    if (!isEndOfRules) rules.push(row);
    else updates.push(row);
  }
  return { rules, updates };
};

const part1 = ({ rules, updates }) => {
  const rulesSet = new Set(rules);
  const validUpdates = [];

  for (const update of updates) {
    const pages = update.split(",");
    let isValidUpdate = true;

    for (let i = 0; i < pages.length - 1; i++) {
      const failingRule = `${pages[i + 1]}|${pages[i]}`;
      if (rulesSet.has(failingRule)) {
        isValidUpdate = false;
        break;
      }
    }
    if (isValidUpdate) validUpdates.push(pages);
  }

  let sum = 0;
  for (const validUpdate of validUpdates)
    sum += Number(validUpdate[Math.floor(validUpdate.length / 2)]);

  console.log(sum);
  return sum;
};

const part2 = ({ rules, updates }) => {
  const correctedUpdates = [];

  for (const update of updates) {
    const rulesSet = new Set();
    const pages = update.split(",").map(Number);

    for (const rule of rules) {
      const [a, b] = rule.split("|").map(Number);
      if (pages.includes(a) && pages.includes(b)) rulesSet.add([a, b]);
    }

    const dependencies = new Map();
    for (const [a, b] of rulesSet)
      dependencies.set(b, (dependencies.get(b) || 0) + 1);

    const corrected = [];
    while (corrected.length < pages.length) {
      for (let i = 0; i < pages.length; i++) {
        const page = pages[i];

        if (corrected.includes(page)) continue;

        if ((dependencies.get(page) || 0) < 1) {
          corrected.push(page);
          for (const [a, b] of rulesSet)
            if (page === a) dependencies.set(b, dependencies.get(b) - 1);
        }
      }
    }

    if (corrected.join(",") !== update) correctedUpdates.push(corrected);
  }

  let sum = 0;
  for (const correctedUpdate of correctedUpdates)
    sum += Number(correctedUpdate[Math.floor(correctedUpdate.length / 2)]);

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
