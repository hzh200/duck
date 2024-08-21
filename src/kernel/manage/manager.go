package manage

import (
	"duck/kernel/constance"
	"duck/kernel/downloaders"
	"duck/kernel/log"
	"duck/kernel/models"
	"duck/kernel/persistence"
	"encoding/json"
	"fmt"
	"time"
)

const MaximumTaskNum = 5

type Manager struct {
	queue *TaskQueue
	persistence *persistence.Persistence
	runningTaskNos []int64
}

func NewManager(persistence *persistence.Persistence) *Manager {
	manager := Manager{}
	manager.queue = NewTaskQueue()
	manager.persistence = persistence
	manager.runningTaskNos = make([]int64, 0)
	return &manager
}

func (manager *Manager) Schedule() {
	go func() {
		for {
			upcomingTaskNos := manager.queue.GetWaitingTaskNos(MaximumTaskNum - len(manager.runningTaskNos))
			for i := 0; i < len(upcomingTaskNos); i++ {
				isContained := false
				for j := 0; j < len(manager.runningTaskNos); j++ {
					if manager.runningTaskNos[j] == upcomingTaskNos[i] {
						isContained = true
						break
					}
				}
				if !isContained {
					manager.runningTaskNos = append(manager.runningTaskNos, upcomingTaskNos[i])
				}
			}
			time.Sleep(500)
		}
	}()
	go func() {
		for {
			for _, runningTaskNo := range manager.runningTaskNos {
				task := manager.queue.GetTaskByNo(runningTaskNo)
				if task.TaskStatus == constance.Waiting {
					task.TaskStatus = constance.Running

					go func() {
						err := downloaders.Download(task)
						if err != nil {
							log.Error(err)
						}
					}()

					go func() {
						for {
							if task.TaskProgress == task.TaskSize {
								task.TaskStatus = constance.Successed
								manager.persistence.UpdateTask(*task)
								break
							}
							manager.persistence.UpdateTask(*task)
							//lint:ignore SA1004 ingore this for now
							time.Sleep(100)
						}
					}()
				}
			} 
			time.Sleep(500)
		}
	}()
}

func (manager *Manager) AddTaskToQueue(task models.Task) {
	manager.persistence.AddTask(&task)
	manager.queue.AddTask(task)
	bytes, _ := json.Marshal(task); log.Info(fmt.Sprintf("Added task: %s", string(bytes)))
}

func (manager Manager) GetAllTasks() []*models.Task {
	return manager.queue.tasks
}

func (manager Manager) GetTaskByNo(taskNo int64) *models.Task {
	return manager.queue.GetTaskByNo(taskNo)
}
