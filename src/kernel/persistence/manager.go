package persistence

import (
	"duck/kernel/models"
	"duck/kernel/persistence/core"
	_ "github.com/mattn/go-sqlite3"
	_ "duck/kernel/persistence/core/dialects"
)

type Manager struct {
	engine *core.Engine
}

func StartManager(dbPath string) error {
	manager := Manager{}
	engine, err := core.New(dbPath, "sqlite3", []interface{}{
		models.Task{}, 
		models.TaskSet{},
	})
	if err != nil {
		return err
	}
	manager.engine = engine
	return nil
}

func (manager *Manager) GetAllTasks() ([]models.Task, error) {
	tasks, err := manager.engine.Select(models.Task{})
	if err != nil {
		return nil, err
	}
	return tasks.([]models.Task), nil
}

func (manager *Manager) AddTask(task models.Task) {
	manager.engine.Insert(task)
}

func (manager *Manager) UpdateTask() {
	
}

func (manager *Manager) RemoveTask() {
	
}

func (manager *Manager) GetAllTaskSets() ([]models.TaskSet, error) {
	taskSets, err := manager.engine.Select(models.TaskSet{})
	if err != nil {
		return nil, err
	}
	return taskSets.([]models.TaskSet), nil
}

func (manager *Manager) AddTaskSet(taskSet models.TaskSet) {
	manager.engine.Insert(taskSet)
}

func (manager *Manager) UpdateTaskSet() {
	
}

func (manager *Manager) RemoveTaskSet() {
	
}
