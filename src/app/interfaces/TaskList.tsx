import React from 'react';
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

import { Task } from '@/models/Task';
import { taskTableColumns, TaskTableColumn } from '@/lib/task-table-columns';
import { cn } from '@/lib/utils';

import { Button } from '@/components/ui/button';

import '@/interfaces/styles/task-list.css';
import { PauseIcon, PlayIcon, ResumeIcon, StopIcon, TrashIcon } from '@radix-ui/react-icons';

interface TaskListProps {
  tasks: Array<Task>
  choosen: Array<number>
  setChoosen: React.Dispatch<React.SetStateAction<number[]>>
}

function TaskList({ tasks, choosen, setChoosen }: TaskListProps) {
  return (
    <div>
      <div id='task-control' className='flex items-center h-10 w-full'>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant='ghost' size='icon' className='h-full w-10'>
              <PlayIcon className='h-6 w-6' />
            </Button>
          </TooltipTrigger>
          <TooltipContent side='top'>Start</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant='ghost' size='icon' className='h-full w-10'>
              <PauseIcon className='h-6 w-6' />
            </Button>
          </TooltipTrigger>
          <TooltipContent side='top'>Pause</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant='ghost' size='icon' className='h-full w-10'>
              <ResumeIcon className='h-6 w-6' />
            </Button>
          </TooltipTrigger>
          <TooltipContent side='top'>Resume</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant='ghost' size='icon' className='h-full w-10'>
              <StopIcon className='h-6 w-6' />
            </Button>
          </TooltipTrigger>
          <TooltipContent side='top'>Stop</TooltipContent>
        </Tooltip>
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant='ghost' size='icon' className='h-full w-10'>
              <TrashIcon className='h-6 w-6' />
            </Button>
          </TooltipTrigger>
          <TooltipContent side='top'>Delete</TooltipContent>
        </Tooltip>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            {taskTableColumns.map((column, index) => (
              <TableHead key={index} className={column.classes ? column.classes.join(' ') : ''}>
                {column.title}
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody className='scroll'>
          {tasks.map((task: Task) => (
            <TableRow 
              key={task.taskNo} 
              className={cn(
                "text-left text-sm transition-all hover:bg-accent",
                choosen.includes(task.taskNo) && "bg-muted"
              )}
              onClick={() => {
                choosen.includes(task.taskNo) ? choosen.splice(choosen.indexOf(task.taskNo), 1) : choosen = [...choosen, task.taskNo];
                setChoosen([...choosen]);
              }}
            >
              {taskTableColumns.map((column: TaskTableColumn, index) => (
                <TableCell key={index} className={column.classes ? column.classes.join(' ') : ''}>
                  {(() => {
                    const value = (task as any)[column.accessorKey];
                    if (value !== (null || undefined)) {
                      return column.convertor ? column.convertor(value, task) : value;
                    }
                    return '';
                  })()}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

export default TaskList;
