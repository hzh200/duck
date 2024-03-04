import { app, BrowserWindow, Event, screen, Tray, Menu, nativeImage } from "electron";
import { resolve, basename } from "node:path";
import { spawn, ChildProcessWithoutNullStreams } from "node:child_process";

const ICON_PATH = resolve(__dirname, "../../asset/favicon.ico");
const EXE_NAME = basename(process.execPath);
const SOFTWARE_NAME = "";

const APP_PATH = resolve(__dirname, "./index.html");
const KERNEL_PATH = resolve(__dirname, "./kernel");
const KERNEL_PORT = app.commandLine.getSwitchValue('port') !== '' && !isNaN(parseInt(app.commandLine.getSwitchValue('port'))) ? parseInt(app.commandLine.getSwitchValue('port')) : 9000;

const [MAIN_WIDTH, MAIN_HEIGHT, DEVTOOLS_WIDTH, DEVTOOLS_HEIGHT] = [1200, 800, 800, 800];

const DEV_MODE: boolean = process.env.NODE_ENV === 'development';
const SILENT_MODE = app.commandLine.getSwitchValue('start-mode') === 'silent';

let tray!: Tray;
let mainWindow!: BrowserWindow;
let devtoolsWindow!: BrowserWindow;
let kernel!: ChildProcessWithoutNullStreams;

const icon = nativeImage.createFromPath(ICON_PATH);

let appTerminating = false;

let screenHeight: number;
let screenWidth: number;
let scaleFactor: number;
// let main_height: number;
// let main_width: number;
// let devtools_height: number;
// let devtools_width: number;

const globalSetting = {
  launchOnStartup: false,
  closeToTray: true
};

if (globalSetting.launchOnStartup) {
  app.setLoginItemSettings({
    openAtLogin: true,
    path: process.execPath,
    args: [
      '--processStart', `"${EXE_NAME}"`,
      '--process-start-args', '"--start-mode=silent"'
    ]
  });
} else {
  app.setLoginItemSettings({ openAtLogin: false });
}

const launchAPP = () => new Promise<void>((resolve, _reject) => {
  const initMainWindow = () => {
    mainWindow = new BrowserWindow({
      width: (MAIN_WIDTH > screenWidth ? screenWidth : MAIN_WIDTH) / scaleFactor,
      height: (MAIN_HEIGHT > screenHeight ? screenHeight : MAIN_HEIGHT) / scaleFactor,
      icon,
      show: false,
      webPreferences: {
        zoomFactor: 1.0 / scaleFactor,
        nodeIntegration: true,
        contextIsolation: false,
      }
    });
  
    mainWindow.loadFile(APP_PATH);
    if (!SILENT_MODE) {
      mainWindow.once("ready-to-show", async () => {
        if (DEV_MODE) {
          await initDevToolsWindow();
        }
        mainWindow.show();
        resolve();
      });
    }
  };

  const initDevToolsWindow = (): Promise<void> => {
    devtoolsWindow = new BrowserWindow({
      width: (MAIN_WIDTH + DEVTOOLS_WIDTH > screenWidth ? screenWidth - MAIN_WIDTH : DEVTOOLS_WIDTH) / scaleFactor,
      height: (DEVTOOLS_HEIGHT > screenHeight ? screenHeight : DEVTOOLS_HEIGHT) / scaleFactor,
      icon,
      show: false,
      webPreferences: {
        zoomFactor: 1.0 / scaleFactor
      }
    });
    devtoolsWindow.setMenu(null);
    mainWindow.webContents.setDevToolsWebContents(devtoolsWindow.webContents);
    mainWindow.webContents.openDevTools({ mode: 'detach' });
  
    const positionDevTools = () => { // Put the DevTools window at the right side of the main window.
      if (devtoolsWindow.isDestroyed()) {
        return;
      }
      const windowBounds = mainWindow.getBounds();
      devtoolsWindow.setPosition(windowBounds.x + windowBounds.width, windowBounds.y);
      devtoolsWindow.setSize(DEVTOOLS_HEIGHT, DEVTOOLS_WIDTH);
    };
  
    const positionMain = () => {
      if (mainWindow.isDestroyed()) {
        return;
      }
      const mainSize = mainWindow.getSize();
      const devtoolSize = devtoolsWindow.isDestroyed() ? [0, 0] : devtoolsWindow.getSize();
      const x = (screenWidth - mainSize[0] - devtoolSize[0]) / 2;
      const y = (screenHeight - Math.max(mainSize[1], devtoolSize[1])) / 2;
      mainWindow.setPosition(Math.floor(x), Math.floor(y));
    };
  
    devtoolsWindow.addListener('closed', positionMain);
    mainWindow.addListener('move', positionDevTools);
    mainWindow.addListener('close', (_event: Event) => {
      if (!devtoolsWindow.isDestroyed()) {
        devtoolsWindow.close();
      }
    });

    positionMain();
    positionDevTools();
  
    return new Promise((resolve) => {
      devtoolsWindow.on('ready-to-show', () => {
        if (!SILENT_MODE) {
          devtoolsWindow.show();
        }
        resolve();
      });
    });
  };

  initMainWindow();
});

const shutdownAPP = () => new Promise<void>((resolve, _reject) => {
  if (!mainWindow.isDestroyed()) {
    mainWindow.close();
  }
  const shutdownPromises = [new Promise<void>((resolve, _reject) => mainWindow.on('closed', resolve))];
  if (DEV_MODE) {
    if (!devtoolsWindow.isDestroyed()) {
      devtoolsWindow.close();
    }
    shutdownPromises.push(new Promise<void>((resolve, _reject) => devtoolsWindow.on('closed', resolve)));
  }
  Promise.all(shutdownPromises).then(() => resolve());
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

const launchTray = () => {
  tray = new Tray(icon);
  const contextMenu = Menu.buildFromTemplate([
    { 
      label: 'Show Main Window', 
      click: mainWindow.show
    },
    { 
      label: 'Quit', 
      click: shutdownSystem 
    }
  ]);
  tray.setToolTip(SOFTWARE_NAME);
  tray.setContextMenu(contextMenu);
  tray.on('double-click', mainWindow.show);
};

const shutdownTray = () => {
  if (!tray.isDestroyed()) {
    tray.destroy();
  }
}

const launchSystem = async () => {
  await launchKernel();
  await launchAPP();
  launchTray();
};

const shutdownSystem = async () => {
  appTerminating = true;
  shutdownTray();
  await shutdownAPP();
  const code: number | null = await shutdownKernel();
  if (code !== null) {
    console.info(`kernel: closed with ${code}`);
  }
  app.quit();
  console.info("System has shutted down.");
  process.exit(1);
};

if (!app.requestSingleInstanceLock()) {
  console.error("An instance is already running.");
  await shutdownSystem();
}

try {
  await app.whenReady();
  ({ height: screenHeight, width: screenWidth } = screen.getPrimaryDisplay().size);
  
  scaleFactor = screen.getPrimaryDisplay().scaleFactor;

  await launchSystem();
  mainWindow.on("close", async (event: Event) => {
    if (appTerminating) {
      return;
    }
    if (globalSetting.closeToTray) {
      mainWindow.hide();
    } else {
      event.preventDefault();
      await shutdownSystem();
    }
  });
} catch (err) {
  console.error(err);
  // await shutdownSystem();
}
