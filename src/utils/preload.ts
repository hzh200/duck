import { contextBridge, ipcRenderer } from 'electron';

contextBridge.exposeInMainWorld('control', {
  minimize: () => ipcRenderer.invoke('minimize'),
  maximize: () => ipcRenderer.invoke('maximize'),
  unmaximize: () => ipcRenderer.invoke('unmaximize'),
  close: () => ipcRenderer.invoke('close')
});

contextBridge.exposeInMainWorld('setting', {
  onSettings: (callback) => ipcRenderer.on('settings', (_event, value) => callback(value)),
  updateSettings: (settingJson) => ipcRenderer.send('updateSettings', settingJson)
});
