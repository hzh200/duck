import * as React from 'react';
import TaskList from '@/interfaces/TaskList';
import TaskInfo from '@/interfaces/TaskInfo';

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';

import { Button } from '@/components/ui/button';
import { ArrowDownIcon, GearIcon, MinusIcon, SquareIcon, Cross2Icon } from '@radix-ui/react-icons';


import './css/app.css';

function AppPage() {
  return (
    <div id='app'>
      <div id='frame' className='h-6'>
        <ResizablePanelGroup
          direction="horizontal"
          className="min-h-[200px] max-w-md rounded-lg border h-6"
        >
          <ResizablePanel defaultSize={80}>
            <div className="flex h-full items-center justify-center p-6">
              <ArrowDownIcon className='h-4 w-4' />
              <span className="font-semibold">Sidebar</span>
            </div>
          </ResizablePanel>
          <ResizableHandle withHandle />
          <ResizablePanel defaultSize={40}>
            <div className="flex h-full items-center justify-center p-6">
              <Button variant='ghost' size='icon' className='h-6 w-6'>
                <GearIcon className='h-4 w-4' />
              </Button>
              <div id='frame-control' className='align-middle h-6'>
                <Button variant='ghost' size='icon' className='h-6 w-6'>
                  <MinusIcon className='h-4 w-4' />
                </Button>
                <Button variant='ghost' size='icon' className='h-6 w-6'>
                  <SquareIcon className='h-4 w-4' />
                </Button>
                <Button variant='ghost' size='icon' className='h-6 w-6'>
                  <Cross2Icon className='h-4 w-4' />
                </Button>
              </div>
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
      <div id='main'>
        <ResizablePanelGroup
          direction="horizontal"
          className="min-h-[200px] max-w-md rounded-lg border"
        >
          <ResizablePanel defaultSize={25}>
            <div className="flex h-full items-center justify-center p-6">
              <TaskList />
            </div>
          </ResizablePanel>
          <ResizableHandle withHandle />
          <ResizablePanel defaultSize={75}>
            <div className="flex h-full items-center justify-center p-6">
              <TaskInfo />
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </div>
    </div>
  );
}

export default AppPage;
