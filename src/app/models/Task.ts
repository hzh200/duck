import TaskStatus from '@/models/TaskStatus';

class Task {
    taskNo: number
    taskName: string
    status: TaskStatus
    progress: number
    size: number
}

export { Task }
