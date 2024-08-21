import React from 'react';
import { Input } from '@/components/ui/input';
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Separator } from '@/components/ui/separator';
import { Label } from "@/components/ui/label";
import { Switch } from '@/components/ui/switch';
import { proxiedSetting } from '@/lib/setting';

function SettingPage(setting: {[key: string]: any}) {
    const onSettingChange = (eleName: string, value: any) => {
        let name = eleName;
        let subName = "";
        if (eleName.includes("-")) {
            [name, subName] = eleName.split("-");
        }
        if (subName === "") {
            proxiedSetting[name] = value;
        } else {
            proxiedSetting[name][subName] = value;
        }
    };

    return (
        <div className='h-full w-full p-4 flex flex-col space-y-4'>
            <div>
                <Label className='text-xl'>Task</Label>
                <Separator className="mt-1 mb-2" />
                <div className="flex flex-col space-y-2">
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='location'>Location:</Label>
                        <div id='location' className="flex items-center space-x-2">
                            <Input name='downloadDirectory' className='w-[200px]' value={setting["downloadDirectory"]} 
                                onChange={event => onSettingChange('downloadDirectory', event.target.value)} />
                        </div>
                    </div>
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='traffic-limit'>Trafic Limit:</Label>
                        <div id='traffic-limit' className="flex items-center space-x-2">
                            <Switch name="trafficLimit-enabled" checked={setting["trafficLimit"]["enabled"]} 
                                onCheckedChange={value => onSettingChange('trafficLimit-enabled', value)} />
                            <Input name="trafficLimit-limit" className='w-[120px]' value={setting["trafficLimit"]["limit"]} disabled={!setting["trafficLimit"]["enabled"]} 
                                onChange={event => onSettingChange('trafficLimit-limit', event.target.value)} />
                            <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                                Mb/s
                            </label>
                        </div>
                    </div>
                </div>
            </div>
            <div>
                <Label className='text-xl'>Network</Label>
                <Separator className="mt-1 mb-2" />
                <div className="flex flex-col space-y-2">
                    <div className="flex items-center space-x-2">
                        <RadioGroup name="proxy-proxyMode" value={setting["proxy"]["proxyMode"]} 
                            onValueChange={value => onSettingChange('proxy-proxyMode', value)}>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="off" id="r1" />
                                <Label htmlFor="r1">Proxy off.</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="system" id="r2" />
                                <Label htmlFor="r2">Use system proxy.</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="manually" id="r3" />
                                <Label htmlFor="r3">Set manually.</Label>
                            </div>
                        </RadioGroup>
                    </div>
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='host'>Host:</Label>
                        <div id='host'>
                            <Input name="proxy-host" className='w-[120px]' value={setting["proxy"]["host"]} disabled={setting["proxy"]["proxyMode"] !== "manually"} 
                                onChange={event => onSettingChange('proxy-host', event.target.value)} />
                        </div>
                        <Label htmlFor='port'>Port:</Label>
                        <div id='port'>
                            <Input name="proxy-port" className='w-[120px]' value={setting["proxy"]["port"]} disabled={setting["proxy"]["proxyMode"] !== "manually"}
                                onChange={event => onSettingChange('proxy-port', event.target.value)} />
                        </div>
                    </div>
                </div>
            </div>     
            <div>
                <Label className='text-xl'>System</Label>
                <Separator className="mt-1 mb-2" />
                <div className="flex flex-col space-y-2">
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='launch-on-startup'>Launch on startup:</Label>
                        <div id="launch-on-startup">
                            <Switch name="launchOnStartup" checked={setting["launchOnStartup"]} 
                                onCheckedChange={value => onSettingChange('launchOnStartup', value)} />
                        </div>
                    </div>
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='slient-mode'>Slient mode:</Label>
                        <div id="slient-mode">
                            <Switch name="slientMode" checked={setting["slientMode"]} 
                                onCheckedChange={value => onSettingChange('slientMode', value)} />
                        </div>
                    </div>
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='close-to-tray'>Close to tray:</Label>
                        <div id="close-to-tray">
                            <Switch name="closeToTray" checked={setting["closeToTray"]} 
                                onCheckedChange={value => onSettingChange('closeToTray', value)} />
                        </div>
                    </div>
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='kernel-port'>Kernel port:</Label>
                        <div id="kernel-port">
                            <Input name="kernelPort" className='w-[120px]' value={setting["kernelPort"]} 
                                onChange={event => {
                                    let port = event.target.value;
                                    let numEndIndex = 0;
                                    for (let i = 0; i < port.length; i++) {
                                        if (isNaN(Number(port.substring(i, i + 1)))) {
                                            break;
                                        }
                                        numEndIndex++;
                                    }
                                    port = port.substring(0, numEndIndex);
                                    onSettingChange('kernelPort', port);
                                }} />
                        </div>
                    </div>
                </div>
            </div>  
        </div>
    )
}

export default SettingPage;
