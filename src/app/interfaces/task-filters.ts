import { Task } from '@/models/Task';
import TaskStatus from "@/models/TaskStatus";

interface TaskFilter {
    name: string
    filter: (tasks: Array<Task>) => Array<Task>
}

const taskFilters: Array<TaskFilter> = [
    {
        name: 'All',
        filter: (tasks: Array<Task>) => tasks.filter(task => task.status)
    },
    {
        name: 'Waiting',
        filter: (tasks: Array<Task>) => tasks.filter(task => task.status === TaskStatus.Waiting)
    },
    {
        name: 'Running',
        filter: (tasks: Array<Task>) => tasks.filter(task => task.status === TaskStatus.Running)
    },
    {
        name: 'Paused',
        filter: (tasks: Array<Task>) => tasks.filter(task => task.status === TaskStatus.Paused)
    },
    {
        name: 'Stopped',
        filter: (tasks: Array<Task>) => tasks.filter(task => task.status === TaskStatus.Stopped)
    },
    {
        name: 'Successed',
        filter: (tasks: Array<Task>) => tasks.filter(task => task.status === TaskStatus.Successed)
    },
    {
        name: 'Failed',
        filter: (tasks: Array<Task>) => tasks.filter(task => task.status === TaskStatus.Failed)
    }
];

export { taskFilters, TaskFilter };
