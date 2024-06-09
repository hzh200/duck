import React from 'react';
import { Input } from '@/components/ui/input';
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Separator } from '@/components/ui/separator';
import { Label } from "@/components/ui/label";
import { Switch } from '@/components/ui/switch';

function SettingPage() {
    return (
        <div className='h-full w-full p-4 flex flex-col space-y-4'>
            <div>
                <Label className='text-xl'>Task</Label>
                <Separator className="mt-1 mb-2" />
                <div className="flex flex-col space-y-2">
                    <div className="flex items-center space-x-2">
                        <Label>Location:</Label>
                        <Input className='w-[200px]' />
                    </div>
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='trafic-limit'>Trafic Limit:</Label>
                        <Switch id="trafic-limit" checked={false} />
                        <Input id='limit' className='w-[120px]' />
                        <label
                            htmlFor="limit"
                            className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                        >
                            Mb/s
                        </label>
                    </div>
                </div>
            </div>
            <div>
                <Label className='text-xl'>Network</Label>
                <Separator className="mt-1 mb-2" />
                <div className="flex flex-col space-y-2">
                    <div className="flex items-center space-x-2">
                        <RadioGroup defaultValue="noProxy">
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="noProxy" id="r1" />
                                <Label htmlFor="r1">No proxy.</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="systemPoxy" id="r2" />
                                <Label htmlFor="r2">Use system proxy.</Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="setManually" id="r3" />
                                <Label htmlFor="r3">Set manually.</Label>
                            </div>
                        </RadioGroup>
                    </div>
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='host'>Host:</Label>
                        <Input id='host' className='w-[120px]' />
                        <Label htmlFor='port'>Port:</Label>
                        <Input id='port' className='w-[120px]' />
                    </div>
                </div>
            </div>     
            <div>
                <Label className='text-xl'>System</Label>
                <Separator className="mt-1 mb-2" />
                <div className="flex flex-col space-y-2">
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='host'>Launch on startup:</Label>
                        <Switch id="host" checked={false} />
                    </div>
                    <div className="flex items-center space-x-2">
                        <Label htmlFor='port'>Close to tray:</Label>
                        <Switch id="port" checked={false} />
                    </div>
                </div>
            </div>  
        </div>
    )
}

export default SettingPage;
