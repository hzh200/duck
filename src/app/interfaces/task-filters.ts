import { Task } from '@/models/Task';
import TaskStatus from "@/models/TaskStatus";

interface TaskFilter {
    name: string
    filter: (tasks: Array<Task>) => Array<Task>
}

const taskFilters: Array<TaskFilter> = [{
    name: 'All',
    filter: (tasks: Array<Task>) => tasks.filter(task => task.status)
}];

const values = Object.values(TaskStatus).filter(v => !isNaN(Number(v)));
const keys = Object.values(TaskStatus).filter(v => isNaN(Number(v)));

for (let i = 0; i < values.length; i++) {
    taskFilters.push({
        name: keys[i].toString(),
        filter: (tasks: Array<Task>) => tasks.filter(task => task.status === values[i])
    });
}

export { taskFilters, TaskFilter };
