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
    }
];

export { taskFilters, TaskFilter };
