import { contextBridge, ipcRenderer } from 'electron';

contextBridge.exposeInMainWorld('control', {
  minimize: () => ipcRenderer.invoke('minimize'),
  maximize: () => ipcRenderer.invoke('maximize'),
  close: () => ipcRenderer.invoke('close')
})
