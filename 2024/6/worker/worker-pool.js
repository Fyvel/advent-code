const { Worker } = require("worker_threads");

module.exports = function CreateWorkerPool(numWorkers, workerScript) {
  const workers = [];
  const queue = [];
  let activeWorkers = 0;
  const resolvers = new Map();

  const handleMessage = (worker, result) => {
    worker.busy = false;
    activeWorkers--;

    const resolver = resolvers.get(worker);
    if (resolver) {
      resolver(result);
      resolvers.delete(worker);
    }

    if (queue.length > 0) {
      const next = queue.shift();
      runTask(next);
    }
  };

  const runTask = (task) => {
    return new Promise((resolve) => {
      const runOnWorker = (t) => {
        const worker = workers.find((w) => !w.busy);
        if (worker) {
          worker.busy = true;
          activeWorkers++;
          resolvers.set(worker, resolve);
          worker.postMessage(t);
        } else {
          setTimeout(() => runOnWorker(t), 0);
        }
      };
      runOnWorker(task);
    });
  };

  const terminate = () => {
    workers.forEach((w) => w.terminate());
  };

  for (let i = 0; i < numWorkers; i++) {
    const worker = new Worker(workerScript);
    worker.on("message", (result) => handleMessage(worker, result));
    workers.push(worker);
  }

  return {
    runTask,
    terminate,
  };
};
