import { app, BrowserWindow, Event } from "electron";
import { resolve } from "node:path";
import { spawn, ChildProcessWithoutNullStreams } from "node:child_process";

const APP_PATH = resolve(__dirname, "./index.html");
const KERNEL_PATH = resolve(__dirname, "./kernel");
const KERNEL_PORT = 9000;

let mainWindow!: BrowserWindow;
let kernel!: ChildProcessWithoutNullStreams;

let appTerminating = false;

const launchAPP = () => new Promise<void>((resolve, _reject) => {
  mainWindow = new BrowserWindow({
    width: 1000,
    height: 800,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false,
    },
    show: false,
  });

  mainWindow.loadFile(APP_PATH);
  mainWindow.once("ready-to-show", () => {
    mainWindow.show();
    resolve();
  });
});

const shutdownAPP = () => new Promise<void>((resolve, _reject) => {
  if (!mainWindow.isDestroyed())
    mainWindow.close();
  mainWindow.on('closed', resolve);
});

const launchKernel = () => new Promise<void>((resolve, _reject) => {
  kernel = spawn(KERNEL_PATH, ["-port", `${KERNEL_PORT}`]);
  kernel.stdout.on("data", (data) => {
    console.log(`kernel: ${data}`);
  });
  kernel.stderr.on("data", (data) => {
    console.error(`kernel: ${data}`);
  });
  resolve();
});

const shutdownKernel = () => new Promise<number | null>((resolve, _reject) => {
  kernel.on("closed", (code) => resolve(code));
  kernel.kill("SIGTERM");
});

const launchSystem = async () => {
  await launchKernel();
  await launchAPP();
};

const shutdownSystem = async () => {
  appTerminating = true;
  await shutdownAPP();
  const code: number | null = await shutdownKernel();
  if (code !== null) {
    console.info(`kernel: closed with ${code}`);
  }
  console.info("System has shutted down.");
  app.quit();
  process.exit(1);
};

if (!app.requestSingleInstanceLock()) {
  console.error("An instance is already running.");
  await shutdownSystem();
}

try {
  await app.whenReady();
  await launchSystem();
  mainWindow.on("close", async (event: Event) => {
    if (appTerminating) {
      return;
    }
    event.preventDefault();
    await shutdownSystem();
  });
} catch (err) {
  console.error(err);
  await shutdownSystem();
}
