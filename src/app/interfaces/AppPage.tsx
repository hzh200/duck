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
import { TooltipProvider } from "@/components/ui/tooltip";

import '@/interfaces/styles/app.css';

import { Task } from '@/models/Task';
import TaskStatus from '@/models/TaskStatus';
import TaskFilterList from './TaskFilterList';
import { taskFilters, TaskFilter } from '../lib/task-filters';

function AppPage() {
  const [layout, setLayout] = useState<number[]>([10, 80, 30]);
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
  const [choosenFilter, setChoosenFilter] = useState<TaskFilter>(taskFilters[0]);
  const [choosenTaskNos, setChoosenTaskNos] = useState<Array<number>>([]);

  return (
    <div id='app' className='h-screen w-full rounded-lg border-solid border-2'>
      <TooltipProvider delayDuration={500}>
        <ResizablePanelGroup
          direction="horizontal"
          className="min-w-full h-full"
          onLayout={(sizes: number[]) => {
            setLayout(sizes);
          }}
        >
          <ResizablePanel defaultSize={layout[0] + layout[1]}>
            <div className='frame h-frame'>
              <div className="flex items-center h-full">
                <Badge variant="outline">
                  <ArrowDownIcon className='h-4 w-4' />
                </Badge>
                <Label>duck</Label>
              </div>
              <div className='drag' />
            </div>
            <ResizablePanelGroup
              direction="horizontal"
              className="min-w-full h-full"
              onLayout={(sizes: number[]) => {
                setLayout(sizes);
              }}
            >
              <ResizablePanel defaultSize={layout[1]} minSize={10} maxSize={10}>
                <div className='main h-main p-2'>
                  <TaskFilterList choosen={choosenFilter} setChoosen={setChoosenFilter} />
                </div>
              </ResizablePanel>
              <ResizableHandle />
              <ResizablePanel defaultSize={layout[1]}>
                <div className='main h-main p-5'>
                  <TaskList 
                    tasks={choosenFilter.filter(tasks)} 
                    choosen={choosenTaskNos.filter(taskNo => choosenFilter.filter(tasks).some(task => task.taskNo === taskNo))} 
                    setChoosen={setChoosenTaskNos} />
                </div>
              </ResizablePanel>
            </ResizablePanelGroup>
          </ResizablePanel>
          <ResizableHandle />
          <ResizablePanel defaultSize={layout[2]} minSize={15} maxSize={30}>
            <div className='frame h-frame'>
              <Button variant='ghost' size='icon' className='h-full w-8'>
                <GearIcon className='h-5 w-5' />
              </Button>
              <div className='drag' />
              <div id='frame-control' className='flex items-center h-full'>
                <Button variant='ghost' size='icon' className='h-full w-8' onClick={(window as any).control.minimize}>
                  <MinusIcon className='h-4 w-4' />
                </Button>
                <Button variant='ghost' size='icon' className='h-full w-8' onClick={(window as any).control.maximize}>
                  <SquareIcon className='h-3 w-3' />
                </Button>
                <Button variant='ghost' size='icon' className='h-full w-8' onClick={(window as any).control.close}>
                  <Cross2Icon className='h-4 w-4' />
                </Button>
              </div>
            </div>
            <div className='main h-main p-5'>
              <TaskInfo task={tasks.find(task => task.taskNo === choosenTaskNos[choosenTaskNos.length - 1])} />
            </div>
          </ResizablePanel>
        </ResizablePanelGroup>
      </TooltipProvider>
    </div>
  );
}

export default AppPage;
