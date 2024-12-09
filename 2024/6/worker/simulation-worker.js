const { parentPort } = require("worker_threads");

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

parentPort.on("message", (data) => {
  const { start, direction, newObstacle, areaLimits, obstacles } = data;
  const areaLimitsSet = new Set(areaLimits);
  const obstaclesSet = new Set(obstacles);
  const result = simulate(
    start,
    direction,
    newObstacle,
    areaLimitsSet,
    obstaclesSet
  );
  parentPort.postMessage(result);
});

const simulate = (start, direction, newObstacle, areaLimits, obstacles) => {
  const localMoves = new Map();
  let pos = start;
  let dir = direction;

  while (true) {
    const key = `${pos[0]}_${pos[1]}`;
    if (localMoves.get(key)?.has(dir)) {
      return true;
    }
    localMoves.set(key, (localMoves.get(key) || new Set()).add(dir));

    const nextPos = [
      pos[0] + positionsMapping[dir][0],
      pos[1] + positionsMapping[dir][1],
    ];
    const nextKey = `${nextPos[0]}_${nextPos[1]}`;

    if (areaLimits.has(nextKey) && nextKey !== newObstacle) {
      return false;
    }

    if (obstacles.has(nextKey) || nextKey === newObstacle) {
      dir = moveMapping[dir];
    } else {
      pos = nextPos;
    }
  }
};
