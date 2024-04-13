import TaskStatus from '@/models/TaskStatus';

class Task {
    taskNo: number
    fileName: string
    status: TaskStatus
    progress: number
    size: number
}

export { Task }
