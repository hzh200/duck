package manage

import (
	"duck/kernel/constance"
	"duck/kernel/models"
)

type TaskQueue struct {
	tasks []*models.Task
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		tasks: make([]*models.Task, 0),
	}
}

func (queue *TaskQueue) AddTask(task models.Task) {
	queue.tasks = append(queue.tasks, &task)
}

func (queue *TaskQueue) GetTaskByNo(taskNo int64) *models.Task {
	for _, task := range queue.tasks {
		if task.TaskNo == taskNo {
			return task
		}
	}
	return nil
}

func (queue *TaskQueue) GetWaitingTaskNos(num int) []int64 {
	count := 0
	res := make([]int64, 0)
	for _, task := range queue.tasks {
		if task.TaskStatus == constance.Waiting {
			res = append(res, task.TaskNo)
			count++
			if count == num {
				break
			}
		}
	}
	return res
}
