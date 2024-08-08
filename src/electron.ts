import { app, BrowserWindow, Event, screen, Tray, Menu, nativeImage, ipcMain } from 'electron';
import path, { resolve, basename } from 'node:path';
import { spawn, ChildProcessWithoutNullStreams } from 'node:child_process';
import { setConsoleMode, Log } from './utils/log';
import os from 'node:os';
import { exit } from 'node:process';

// Running environment.
let platform = ""
if (os.type() == 'Windows_NT') {
  platform = "windows";
} else if (os.type() == 'Linux') {
  platform = "linux";
} else {
  exit();
}

// External resource paths.
const ICON_PATH = resolve(__dirname, '../../asset/favicon.ico');
const EXE_NAME = basename(process.execPath);
const SOFTWARE_NAME = '';
const APP_PATH = resolve(__dirname, './index.html');
const KERNEL_PATH = resolve(__dirname, './kernel' + (platform === 'windows' ? '.exe' : ''));
const KERNEL_PORT = app.commandLine.getSwitchValue('port') !== '' && !isNaN(parseInt(app.commandLine.getSwitchValue('port'))) ? parseInt(app.commandLine.getSwitchValue('port')) : 9000;
const DB_PATH = resolve(__dirname, './duck.db')

// System config.
const DEV_MODE: boolean = process.env.NODE_ENV === 'debug';
const SILENT_MODE = app.commandLine.getSwitchValue('start-mode') === 'silent';

// System setting.
const globalSetting = {
  launchOnStartup: false,
  closeToTray: true
};

// System resources.
let tray!: Tray;
let mainWindow!: BrowserWindow;
let devtoolsWindow!: BrowserWindow;
let kernel!: ChildProcessWithoutNullStreams;
const icon = nativeImage.createFromPath(ICON_PATH);

let appTerminating = false;

// Interface parameters.
let [screenWidth, screenHeight]: [number, number] = [0, 0];
let [mainWidth, mainHeight]: [number, number] = [0, 0];
let [devtoolsWidth, devtoolsHeight]: [number, number] = [0, 0];
let [mainX, mainY]: [number, number] = [0, 0];
let [devtoolsX, devtoolsY]: [number, number] = [0, 0];
let [xBeforeMaximization, yBeforeMaximization]: [number, number] = [0, 0];
let [widthBeforeMaximization, heightBeforeMaximization]: [number, number] = [0, 0];
let scaleFactor: number;

const positionMain = (x: number, y: number) => {
  if (!mainWindow || mainWindow.isDestroyed()) {
    return;
  }
  mainWindow.setPosition(x, y);
};

const positionDevTools = (x: number, y: number) => {
  if (!devtoolsWindow || devtoolsWindow.isDestroyed()) {
    return;
  }
  devtoolsWindow.setPosition(x, y);
};

const resizeMain = (width: number, height: number) => {
  if (!mainWindow || mainWindow.isDestroyed()) {
    return;
  }
  mainWindow.setSize(width, height);
};

const launchAPP = () => new Promise<void>((resolve, _reject) => {
  const initMainWindow = () => {
    mainWindow = new BrowserWindow({
      width: mainWidth / scaleFactor,
      height: mainHeight / scaleFactor,
      icon,
      show: false,
      frame: false,
      transparent: true,
      webPreferences: {
        zoomFactor: 1.0 / scaleFactor,
        preload: path.resolve(__dirname, 'preload.js')
      }
    });
  
    mainWindow.setMenu(null);
    mainWindow.loadFile(APP_PATH);

    if (!SILENT_MODE) {
      mainWindow.once("ready-to-show", async () => {
        mainWindow.show();
        if (DEV_MODE) {
          await initDevToolsWindow();
        }
        resolve();
      });
    }

    mainWindow.addListener("move", () => {
      ({x: mainX, y: mainY} = mainWindow.getBounds());
    });

    mainWindow.addListener("resize", () => {
      ({width: mainWidth, height: mainHeight} = mainWindow.getBounds());
    });
  };

  const initDevToolsWindow = (): Promise<void> => {
    devtoolsWindow = new BrowserWindow({
      width: devtoolsWidth / scaleFactor,
      height: devtoolsHeight / scaleFactor,
      icon,
      show: false,
      frame: true,
      webPreferences: {
        zoomFactor: 1.0 / scaleFactor
      }
    });

    devtoolsWindow.setMenu(null);
    mainWindow.webContents.setDevToolsWebContents(devtoolsWindow.webContents);
    mainWindow.webContents.openDevTools({ mode: 'detach' });
 
    devtoolsWindow.addListener('closed', () => {
      mainX = Math.floor((screenWidth - mainWidth) / 2);
      mainY = Math.floor((screenHeight - mainHeight) / 2);
      positionMain(mainX, mainY);
    });

    mainWindow.addListener('move', () => {
      devtoolsX = mainX + mainWidth;
      devtoolsY = mainY;
      positionDevTools(devtoolsX, devtoolsY);
    });

    mainWindow.addListener("resize", () => {
      devtoolsX = mainX + mainWidth;
      devtoolsY = mainY;
      positionDevTools(devtoolsX, devtoolsY);
    });

    mainWindow.addListener('close', (_event: Event) => {
      if (!devtoolsWindow.isDestroyed()) {
        devtoolsWindow.close();
      }
    });

    // mainWindow maximize and minimize event can't be triggered here
    // mainWindow.on('minimize', (_event: Event) => {});
    // mainWindow.on('maximize', (_event: Event) => {});
    // mainWindow.on('unmaximize', (_event: Event) => {});

    positionMain(mainX, mainY);
    positionDevTools(devtoolsX, devtoolsY);
  
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
  kernel = spawn(KERNEL_PATH, [
    '-mode', 
    `${DEV_MODE ? 'dev' : 'proc'}`, 
    '-port', 
    `${KERNEL_PORT}`, 
    '-dbPath',
    `${DB_PATH}`
  ]);
  kernel.stdout.on("data", (data) => {
    Log.kernelInfo(data);
  });
  kernel.stderr.on("data", (data) => {
    Log.kernelError(data);
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
    Log.info(`kernel closed with ${code}`);
  }
  
  Log.info("System has been shutted down.");
  app.quit();
  process.exit(1);
};

if (!app.requestSingleInstanceLock()) {
  Log.error("An instance is already running.");
  await shutdownSystem();
}

try {
  await app.whenReady();
  setConsoleMode(DEV_MODE);
  
  const screenInfo = screen.getPrimaryDisplay();

  ({ width: screenWidth, height: screenHeight } = screenInfo.size);
  mainWidth = DEV_MODE ? Math.ceil((screenWidth * 0.98) * 0.6) : Math.ceil(screenWidth * 0.8);
  mainHeight = Math.ceil(screenHeight * 0.6);
  devtoolsWidth = Math.ceil((screenWidth * 0.98) * 0.4);
  devtoolsHeight = Math.ceil(screenHeight * 0.8);
  // Defaultly, put mainWindow on the center of the screen, if in dev mode, put mainWindow and devWindow together on the center of the screen.
  mainX = Math.floor((screenWidth - mainWidth) / 2);
  if (DEV_MODE) {
    mainX = Math.floor(mainX - devtoolsWidth / 2);
  }
  mainY = Math.floor((screenHeight - (DEV_MODE ? Math.max(mainHeight, devtoolsHeight) : mainHeight)) / 2);
  // Put the DevTools window at the right side of the main window.
  devtoolsX = mainX + mainWidth;
  devtoolsY = mainY;

  scaleFactor = screenInfo.scaleFactor;

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

  ipcMain.handle('minimize', () => {
    if (mainWindow && !mainWindow.isDestroyed() && mainWindow.minimizable) {
      mainWindow.minimize();
      // if (devtoolsWindow && !devtoolsWindow.isDestroyed()) {
      //   // Not working.
      //   mainWindow.addListener("restore", devtoolsWindow.show);
      //   devtoolsWindow.hide();
      //   // Would cause error.
      //   devtoolsWindow.minimize();
      // }
    }
  });
  ipcMain.handle('maximize', () => {
    // if (mainWindow && !mainWindow.isDestroyed() && mainWindow.maximizable) {
    if (mainWindow && !mainWindow.isDestroyed()) {
      // mainWindow.maximize();
      [xBeforeMaximization, yBeforeMaximization] = [mainX, mainY];
      [widthBeforeMaximization, heightBeforeMaximization] = [mainWidth, mainHeight];
      positionMain(0, 0);
      resizeMain(screenWidth, screenHeight);
      if (devtoolsWindow && !devtoolsWindow.isDestroyed() && devtoolsWindow.isVisible()) {
        devtoolsWindow.hide();
      }
    }
  });
  ipcMain.handle('unmaximize', () => {
    // if (mainWindow && !mainWindow.isDestroyed() && mainWindow.isMaximized()) {
    if (mainWindow && !mainWindow.isDestroyed()) {
      // mainWindow.unmaximize();
      positionMain(xBeforeMaximization, yBeforeMaximization);
      resizeMain(widthBeforeMaximization, heightBeforeMaximization);
      if (devtoolsWindow && !devtoolsWindow.isDestroyed() && !devtoolsWindow.isVisible()) {
        devtoolsWindow.show();
      }
    }
  });
  ipcMain.handle('close', () => {
    if (!mainWindow.isDestroyed()) {
      mainWindow.close();
    }
  });
} catch (err) {
  Log.error(err instanceof Error ? err : String(err));
  await shutdownSystem();
}
