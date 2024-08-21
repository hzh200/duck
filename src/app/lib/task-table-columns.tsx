import React from 'react';

import { Task } from '@/models/Task';
import TaskStatus from "@/models/TaskStatus";
import { Progress } from "@/components/ui/progress";
import { CheckIcon, ClockIcon, Cross2Icon } from "@radix-ui/react-icons";

interface TaskTableColumn {
    title: string,
    accessorKey: string,
    classes?: Array<string>
    convertor?: (value :any, task: Task) => any,
}

const taskTableColumns: Array<TaskTableColumn> = [
    {
        'title': 'FileName',
        'accessorKey': 'taskName',
        'classes': ['w-10']
    },
    {
        'title': 'Status',
        'accessorKey': 'status',
        'classes': ['w-10'],
        'convertor': (status: TaskStatus, task: Task) => 
            status === TaskStatus.Successed ? <CheckIcon className='h-4 w-4' /> : 
            status === TaskStatus.Stopped ? <CheckIcon className='h-4 w-4' /> :
            status === TaskStatus.Running ? <Progress value={
                'progress' in task && 'size' in task ? Math.floor(task['progress'] * 100 / task['size']) : 0
            } /> :
            status === TaskStatus.Waiting ? <ClockIcon className='h-4 w-4' /> :
            status === TaskStatus.Failed ? <Cross2Icon className='h-4 w-4' /> : <span /> 
    }
];

export { taskTableColumns, TaskTableColumn };
