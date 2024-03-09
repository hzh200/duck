import React, { useState, useEffect } from 'react';
import TaskList from '@/interfaces/TaskList';
import TaskInfo from '@/interfaces/TaskInfo';

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable';

import { Button } from '@/components/ui/button';
import { ArrowDownIcon, GearIcon, MinusIcon, SquareIcon, Cross2Icon } from '@radix-ui/react-icons';
import { Badge } from "@/components/ui/badge";
import { Label } from "@/components/ui/label";

import '@/interfaces/styles/app.css';

import { Task } from '@/models/Task';
import TaskStatus from '@/models/TaskStatus';

function AppPage() {
  const [layout, setLayout] = useState<number[]>([20, 180, 80]);
  const [tasks, setTasks] = useState<Array<Task>>([{
    'taskNo': 1,
    'fileName': 'a',
    'status': TaskStatus.Successed,
    'progress': 1,
    'size': 1
  }, {
    'taskNo': 2,
    'fileName': 'b',
    'status': TaskStatus.Running,
    'progress': 66,
    'size': 100
  }]);

  return (
    <div id='app' className='h-screen w-full rounded-lg border-solid border-2'>
      <ResizablePanelGroup
        direction="horizontal"
        className="min-w-full h-full"
        onLayout={(sizes: number[]) => {
          setLayout(sizes);
        }}
      >
        <ResizablePanel 
          defaultSize={layout[0]}
          minSize={15}
          maxSize={20}
        >
          <div className='frame h-frame'>
            <div className="flex items-center h-full">
              <Badge variant="outline">
                <ArrowDownIcon className='h-4 w-4' />
              </Badge>
              <Label>duck</Label>
            </div>
          </div>
          <div className='main h-main'>

          </div>
        </ResizablePanel>
        <ResizableHandle />
        <ResizablePanel 
          defaultSize={layout[1]}
        >
          <div className='frame h-frame' />
          <div className='main h-main'>
            <TaskList tasks={tasks} />
          </div>
        </ResizablePanel>
        <ResizableHandle />
        <ResizablePanel 
          defaultSize={layout[2]}
          minSize={15}
        >
          <div className='frame h-frame'>
            <div className="flex items-center justify-between h-full">
              <Button variant='ghost' size='icon' className='h-full w-6'>
                <GearIcon className='h-4 w-4' />
              </Button>
              <div id='frame-control' className='flex items-center h-full'>
                <Button variant='ghost' size='icon' className='h-full w-8'>
                  <MinusIcon className='h-4 w-4' />
                </Button>
                <Button variant='ghost' size='icon' className='h-full w-8'>
                  <SquareIcon className='h-3 w-3' />
                </Button>
                <Button variant='ghost' size='icon' className='h-full w-8'>
                  <Cross2Icon className='h-4 w-4' />
                </Button>
              </div>
            </div>
          </div>
          <div className='main h-main'>
            <TaskInfo />
          </div>
        </ResizablePanel>
      </ResizablePanelGroup>
    </div>
  );
}

export default AppPage;
