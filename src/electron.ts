import { app, BrowserWindow, Event, screen, Tray, Menu, nativeImage, ipcMain } from 'electron';
import path, { resolve, basename } from 'node:path';
import { spawn, ChildProcessWithoutNullStreams } from 'node:child_process';
import { setConsoleMode, Log } from './utils/log';

const ICON_PATH = resolve(__dirname, '../../asset/favicon.ico');
const EXE_NAME = basename(process.execPath);
const SOFTWARE_NAME = '';

const APP_PATH = resolve(__dirname, './index.html');
const KERNEL_PATH = resolve(__dirname, './kernel');
const KERNEL_PORT = app.commandLine.getSwitchValue('port') !== '' && !isNaN(parseInt(app.commandLine.getSwitchValue('port'))) ? parseInt(app.commandLine.getSwitchValue('port')) : 9000;

const DB_PATH = resolve(__dirname, './duck.db')

const DEV_MODE: boolean = process.env.NODE_ENV === 'development';
const SILENT_MODE = app.commandLine.getSwitchValue('start-mode') === 'silent';

let tray!: Tray;
let mainWindow!: BrowserWindow;
let devtoolsWindow!: BrowserWindow;
let kernel!: ChildProcessWithoutNullStreams;

const icon = nativeImage.createFromPath(ICON_PATH);

let appTerminating = false;

let screenWidth: number;
let screenHeight: number;
let mainHeight: number;
let mainWidth: number;
let devtoolsHeight: number;
let devtoolsWidth: number;
let scaleFactor: number;

const globalSetting = {
  launchOnStartup: false,
  closeToTray: true
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
        if (DEV_MODE) {
          await initDevToolsWindow();
        }
        mainWindow.show();
        resolve();
      });
    }
  };

  const positionDevTools = () => { // Put the DevTools window at the right side of the main window.
    if (devtoolsWindow.isDestroyed()) {
      return;
    }
    const windowBounds = mainWindow.getBounds();
    devtoolsWindow.setPosition(windowBounds.x + windowBounds.width, windowBounds.y);
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
 
    devtoolsWindow.addListener('closed', positionMain);
    mainWindow.addListener('move', positionDevTools);
    mainWindow.addListener('close', (_event: Event) => {
      if (!devtoolsWindow.isDestroyed()) {
        devtoolsWindow.close();
      }
    });
    mainWindow.addListener('minimize', (_event: Event) => {
      if (devtoolsWindow && !devtoolsWindow.isDestroyed() && devtoolsWindow.minimizable) {
        devtoolsWindow.minimize();
      }
    });
    // mainWindow.addListener('maximize', (_event: Event) => {
    //   if (!devtoolsWindow.isDestroyed() && devtoolsWindow.maximizable) {
    //     devtoolsWindow.maximize();
    //   }
    // });

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
  
  ({ height: screenHeight, width: screenWidth } = screen.getPrimaryDisplay().size);
  mainHeight = Math.ceil(screenHeight * 0.6);
  devtoolsHeight = Math.ceil(screenHeight * 0.8);
  mainWidth = DEV_MODE ? Math.ceil((screenWidth * 0.98) * 0.6) : Math.ceil(screenWidth * 0.8);
  devtoolsWidth = Math.ceil((screenWidth * 0.98) * 0.4);
  scaleFactor = screen.getPrimaryDisplay().scaleFactor;

  
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
    if (mainWindow && mainWindow.minimizable) {
      mainWindow.minimize();
    }
  });
  ipcMain.handle('maximize', () => {
    if (mainWindow && mainWindow.maximizable) {
      mainWindow.maximize();
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
