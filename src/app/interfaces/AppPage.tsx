import React, { useEffect, useState } from 'react';

import { Button } from '@/components/ui/button';
import { ArrowDownIcon, GearIcon, MinusIcon, SquareIcon, Cross2Icon, CopyIcon, PlusIcon, RowsIcon, HamburgerMenuIcon, ListBulletIcon } from '@radix-ui/react-icons';
import { Badge } from "@/components/ui/badge";
import { Label } from "@/components/ui/label";
import { Toaster } from "@/components/ui/sonner"

import '@/interfaces/styles/app.css';
import Router from '@/router';
import { proxiedSetting } from '@/lib/setting';

function AppPage() {
  const [isWindowMaximized, setIsWindowMaximized] = useState<boolean>(false);
  const [path, setPath] = useState<string>('/task');
  const [setting, setSetting] = useState<{[key: string]: any}>({});

  useEffect(() => {
    setInterval(() => setSetting({...proxiedSetting}), 10);
  }, []);

  return (
    <React.Fragment>
      <div id='app' className='h-screen w-full rounded-lg border-solid border'>
        <div className='flex h-frame border-solid border-b'>
          <div className="flex items-center h-full">
            <Badge variant="outline" className='mx-2 my-1'>
              <ArrowDownIcon className='h-6 w-4' />
            </Badge>
            <Label className='text-lg'>duck</Label>
          </div>
          <div className='drag' />
          <div id='frame-control' className='flex items-center h-full'>
            <Button variant='ghost' size='icon' className='h-full w-12' onClick={(window as any).control.minimize}>
              <MinusIcon className='h-4 w-4' />
            </Button>
            <Button variant='ghost' size='icon' className='h-full w-12' onClick={() => {
              if (isWindowMaximized) {
                (window as any).control.unmaximize()
              } else {
                (window as any).control.maximize()
              }
              setIsWindowMaximized(!isWindowMaximized);
            }}>
              { isWindowMaximized ? <CopyIcon className='transform -scale-x-100 h-3 w-3' /> : <SquareIcon className='h-3 w-3' /> }
            </Button>
            <Button variant='ghost' size='icon' className='h-full w-12' onClick={(window as any).control.close}>
              <Cross2Icon className='h-4 w-4' />
            </Button>
          </div>
        </div>
        <div className='flex h-body w-full'>
          <nav className='flex flex-col justify-between h-full w-nav border-solid border-r rounded-none pb-1'>
            <div className='flex flex-col w-full'>
              <Button variant='ghost' size='icon' className={`h-12 w-full rounded-none ${path === '/download' ? 'choosen' : ''}`}
                onClick={() => setPath('/download')} >
                <PlusIcon className='h-8 w-8'/>
              </Button>
              <Button variant='ghost' size='icon' className={`h-12 w-full rounded-none ${path === '/task' ? 'choosen' : ''}`}
                onClick={() => setPath('/task')} >
                <ListBulletIcon className='h-8 w-8'/>
              </Button>
            </div>
            <div className='w-full'>
              <Button variant='ghost' size='icon' className={`h-12 w-full rounded-none ${path === '/setting' ? 'choosen' : ''}`}
                onClick={() => setPath('/setting')} >
                <GearIcon className='h-8 w-8'/>
              </Button>
            </div>
          </nav>
          <main className='h-full w-main'>
            {Router.route(path, setting)}
          </main>
        </div>
      </div>
      <Toaster />
    </React.Fragment>
    
  );
}

export default AppPage;
