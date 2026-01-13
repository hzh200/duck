import { app, BrowserWindow, Event, screen, Tray, Menu, nativeImage, ipcMain, webFrame } from 'electron';
import path, { resolve, basename } from 'node:path';
import fs from 'node:fs';
import { spawn, ChildProcessWithoutNullStreams } from 'node:child_process';
import { setConsoleMode, Log } from './electron/log';
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
const DB_PATH = resolve(__dirname, './duck.db')
const SETTING_PATH = resolve(__dirname, './setting.json');

// System config.
const IS_DEV_MODE: boolean = process.env.NODE_ENV === 'development';

// System setting.
const globalSetting = {
  launchOnStartup: false,
  slientMode: false,
  closeToTray: true,
  kernelPort: 9000,
  downloadDirectory: path.join(os.homedir(), "Downloads"),
  proxy: {
    proxyMode: "system", // off | system | manually
    host: "",
    port: ""
  },
  trafficLimit: {
    enabled: false,
    limit: 0
  }
};

const kernelConfig = {
  "-mode": IS_DEV_MODE ? "dev" : "proc",
  "-dbPath": DB_PATH,
  "-settingPath": SETTING_PATH
}

// System resources.
let tray!: Tray;
let mainWindow!: BrowserWindow;
let devtoolsWindow!: BrowserWindow;
let settingTimeout!: NodeJS.Timeout;
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

//// Position ////

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
  mainWindow.setSize(width, height, true);
};

const launchAPP = () => new Promise<void>((resolve, _reject) => {
  const initMainWindow = () => {
    mainWindow = new BrowserWindow({
      // width: mainWidth / scaleFactor,
      // height: mainHeight / scaleFactor,
      width: mainWidth,
      height: mainHeight,
      icon,
      show: false,
      frame: true,
      // transparent: true,
      webPreferences: {
        preload: path.resolve(__dirname, 'preload.js'),
        zoomFactor: 1.0 / scaleFactor
      }
    });
  
    // mainWindow.setMenu(null);
    mainWindow.loadFile(APP_PATH);

    if (!globalSetting.slientMode) {
      mainWindow.once("ready-to-show", async () => {
        if (IS_DEV_MODE) {
          await initDevToolsWindow();
        }
        mainWindow.show();
        settingTimeout = setInterval(() => mainWindow.webContents.send('settings', globalSetting), 100);
        resolve();
      });
    }

    // Showing mainWindow would also trigger this event.
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

    // devtoolsWindow would automatically increase size while holding mainWindow under 'move' event.
    mainWindow.addListener('moved', () => {
      devtoolsX = mainX + mainWidth;
      devtoolsY = mainY;
      positionDevTools(devtoolsX, devtoolsY);
    });

    mainWindow.addListener("resized", () => {
      devtoolsX = mainX + mainWidth;
      devtoolsY = mainY;
      positionDevTools(devtoolsX, devtoolsY);
    });

    // mainWindow maximize and minimize event wouldn't be triggered here
    // mainWindow.on('minimize', (_event: Event) => {});
    // mainWindow.on('maximize', (_event: Event) => {});
    // mainWindow.on('unmaximize', (_event: Event) => {});

    positionMain(mainX, mainY);
    positionDevTools(devtoolsX, devtoolsY);
  
    return new Promise((resolve) => {
      devtoolsWindow.on('ready-to-show', () => {
        if (!globalSetting.slientMode) {
          devtoolsWindow.show();
        }
        resolve();
      });
    });
  };

  initMainWindow();
});

const shutdownAPP = () => {
  clearInterval(settingTimeout);
  if (!mainWindow.isDestroyed()) {
    mainWindow.destroy();
  }
  if (IS_DEV_MODE) {
    if (!devtoolsWindow.isDestroyed()) {
      devtoolsWindow.destroy();
    }
  }
};

const launchKernel = () => new Promise<void>((resolve, _reject) => {
  kernel = spawn(KERNEL_PATH, Object.entries(kernelConfig).flat());
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
      click: () => mainWindow.show()
    },
    { 
      label: 'Quit', 
      click: shutdownSystem 
    }
  ]);
  tray.setToolTip(SOFTWARE_NAME);
  tray.setContextMenu(contextMenu);
  tray.on('double-click', () => {
    mainWindow.show();
    devtoolsWindow.show();
  });
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
  Log.info("the system has been started.");
};

const shutdownSystem = async () => {
  appTerminating = true;
  shutdownTray();
  shutdownAPP();
  const code: number | null = await shutdownKernel();
  if (code !== null) {
    Log.info(`kernel closed with ${code}`);
  }
  
  Log.info("the system has been shutted down.");
  
  await new Promise<void>((resolve, reject) => {
    setTimeout(() => resolve(), 10000)
  })
  app.quit();
  process.exit(1);

  // process.on('exit', () => process.exit(1));
};

const initFrame = () => {
    const screenInfo = screen.getPrimaryDisplay();

    ({ width: screenWidth, height: screenHeight } = screenInfo.size);
    scaleFactor = screenInfo.scaleFactor;

    mainWidth = IS_DEV_MODE ? Math.ceil((screenWidth * 0.85) * 0.6) : Math.ceil(screenWidth * 0.7);
    mainHeight = Math.ceil(screenHeight * 0.6);
    devtoolsWidth = Math.ceil((screenWidth * 0.8) * 0.4);
    devtoolsHeight = Math.ceil(screenHeight * 0.8);
    // By default, put mainWindow on the center of the screen, if in dev mode, put mainWindow and devWindow together on the center of the screen.
    mainX = Math.ceil((screenWidth - mainWidth) / 2);
    if (IS_DEV_MODE) {
      mainX = Math.ceil(mainX - devtoolsWidth / 2);
    }
    mainY = Math.ceil((screenHeight - (IS_DEV_MODE ? Math.max(mainHeight, devtoolsHeight) : mainHeight)) / 2);
    // Put the DevTools window at the right side of the main window.
    devtoolsX = mainX + mainWidth;
    devtoolsY = mainY;
};

// Listeners for user-controlled panels.
const addUserControllListeners = () => {
  ipcMain.handle('minimize', () => {
    if (mainWindow && !mainWindow.isDestroyed() && mainWindow.minimizable) {
      mainWindow.minimize();
    }
  });
  ipcMain.handle('maximize', () => {
    if (mainWindow && !mainWindow.isDestroyed()) {
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
    if (mainWindow && !mainWindow.isDestroyed()) {
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
};

const main = async () => {
  if (!app.requestSingleInstanceLock()) {
    Log.error("an instance is already running.");
    await shutdownSystem();
  }

  try {
    await app.whenReady();
    setConsoleMode(IS_DEV_MODE);

    initFrame();

    if (fs.existsSync(SETTING_PATH)) {
      Object.assign(globalSetting, JSON.parse(fs.readFileSync(SETTING_PATH, {encoding: "utf-8"}))) 
    } else {
      fs.writeFileSync(SETTING_PATH, JSON.stringify(globalSetting))
    }

    ipcMain.on('updateSettings', (_event, settingsJson) => {
      Object.assign(globalSetting, JSON.parse(settingsJson));
      fs.writeFileSync(SETTING_PATH, settingsJson)
    });

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
      event.preventDefault();
      if (globalSetting.closeToTray) {
        mainWindow.hide();
        devtoolsWindow.hide();
      } else {
        await shutdownSystem();
      }
    });

    // addUserControllListeners();
  } catch (err) {
    Log.error(err instanceof Error ? err : String(err));
    await shutdownSystem();
  }
};

// Main process.
if (process.type === "browser") {
  main();
}
