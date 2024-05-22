import { contextBridge, ipcRenderer } from 'electron';

contextBridge.exposeInMainWorld('control', {
  minimize: () => ipcRenderer.invoke('minimize'),
  maximize: () => ipcRenderer.invoke('maximize'),
  unmaximize: () => ipcRenderer.invoke('unmaximize'),
  close: () => ipcRenderer.invoke('close')
})
