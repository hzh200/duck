let globalSetting = {};
let proxiedSetting = {};

(<any>window).setting.onSettings((setting: { [key: string]: any }) => {
    globalSetting = setting;
    proxiedSetting = setSettingProxyHandler(setting);
});

const settingProxyHandler: ProxyHandler<any> = {
    set(target: any, prop: string | symbol, value: any, _receiver: any): boolean {
        target[prop] = value;
        for (let p in globalSetting) {
            globalSetting[p] = proxiedSetting[p];
        }
        (<any>window).setting.updateSettings(JSON.stringify(globalSetting));
        return true;
    }
};

const setSettingProxyHandler = (setting) => {
    setSettingProxyHandlerCore(setting);
    setting = new Proxy(setting, settingProxyHandler);
    return setting;
};

const setSettingProxyHandlerCore = (obj: { [key: string]: any }): void => {
    for (let prop in obj) {
        if (typeof obj[prop] === 'object' && obj[prop]) {
            setSettingProxyHandlerCore(obj[prop]);
            obj[prop] = new Proxy(obj[prop], settingProxyHandler);
        }
    }
};

export { proxiedSetting };
