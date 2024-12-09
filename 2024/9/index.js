const fs = require("fs").promises;
const path = require("path");

const readData = async () => {
  const data = await fs.readFile(path.join(__dirname, "data.txt"), "utf-8");
  return data.split("\n");
};

const formatData = async (rows) => rows[0];

/** @param {string} diskMap */
const part1 = (diskMap) => {
  const fileBlocksArray = [];
  let blockId = 0;

  for (let i = 0; i < diskMap.length; i++) {
    const count = parseInt(diskMap[i]);
    for (let j = 0; j < count; j++) {
      fileBlocksArray.push(i % 2 === 0 ? blockId.toString() : ".");
    }
    if (i % 2 === 0) {
      blockId++;
    }
  }

  let leftIndex = 0;
  let rightIndex = fileBlocksArray.length - 1;

  while (leftIndex < rightIndex) {
    const left = fileBlocksArray[leftIndex];
    const right = fileBlocksArray[rightIndex];

    if (left === "." && right !== ".") {
      const temp = fileBlocksArray[leftIndex];
      fileBlocksArray[leftIndex] = fileBlocksArray[rightIndex];
      fileBlocksArray[rightIndex] = temp;

      leftIndex++;
      rightIndex--;
      continue;
    }

    if (right === ".") {
      rightIndex--;
    }
    if (left !== ".") {
      leftIndex++;
    }
  }

  let checksum = 0n;
  for (let i = 0; i < fileBlocksArray.length; i++) {
    if (fileBlocksArray[i] === ".") break;
    checksum += BigInt(i) * BigInt(fileBlocksArray[i]);
  }

  // console.log(fileBlocksArray.join(""));
  console.log(checksum.toString());
  return checksum;
};

const createBlock = (value, start, size) => ({ value, start, size });

const slidingWindow = (fileBlocks) => {
  const result = [...fileBlocks];
  const blocks = [];
  let currentValue = result[0];
  let currentStart = 0;
  let currentSize = 1;

  for (let i = 1; i < result.length; i++) {
    if (result[i] === currentValue) {
      currentSize++;
    } else {
      blocks.push(createBlock(currentValue, currentStart, currentSize));
      currentValue = result[i];
      currentStart = i;
      currentSize = 1;
    }
  }
  blocks.push(createBlock(currentValue, currentStart, currentSize));

  // right to left
  for (let i = blocks.length - 1; i >= 0; i--) {
    const block = blocks[i];
    if (block.value === ".") continue;

    let freeSpaceIndex = -1;
    /** @type {ReturnType<typeof createBlock>} */
    let leftMostFreeSpace;

    for (let j = 0; j < i; j++) {
      const localBlock = blocks[j];
      if (localBlock.value !== ".") continue;
      if (localBlock.size >= block.size) {
        leftMostFreeSpace = localBlock;
        freeSpaceIndex = j;
        break;
      }
    }

    if (freeSpaceIndex !== -1) {
      for (let k = 0; k < block.size; k++) {
        result[leftMostFreeSpace.start + k] = result[block.start + k];
      }

      for (let k = block.start; k < block.start + block.size; k++) {
        result[k] = ".";
      }

      blocks[freeSpaceIndex].size -= block.size;
      blocks[freeSpaceIndex].start += block.size;
    }
  }

  return result;
};

/** @param {string} diskMap */
const part2 = (diskMap) => {
  const fileBlocksArray = [];
  let blockId = 0;

  for (let i = 0; i < diskMap.length; i++) {
    const count = parseInt(diskMap[i]);
    for (let j = 0; j < count; j++) {
      fileBlocksArray.push(i % 2 === 0 ? blockId.toString() : ".");
    }
    if (i % 2 === 0) {
      blockId++;
    }
  }
  // console.log(fileBlocksArray.join(""));

  const result = slidingWindow(fileBlocksArray);
  // console.log(result.join(""));

  let checksum = 0n;
  for (let i = 0; i < result.length; i++) {
    if (result[i] === ".") continue;
    checksum += BigInt(i) * BigInt(result[i]);
  }

  console.log(checksum.toString());
  return checksum;
};

readData()
  .then(formatData)
  .then(async (formattedData) => {
    part1(formattedData);
    part2(formattedData);
  })
  .catch((err) => console.error("Get rekt:", err));
