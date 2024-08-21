package persistence

import (
	"duck/kernel/models"
	"duck/kernel/persistence/core"
	_ "duck/kernel/persistence/core/dialects"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Persistence struct {
	engine *core.Engine
}

func InitPersistence(dbPath string) (*Persistence, error) {
	persistence := Persistence{}
	
	engine, err := core.New(dbPath, "sqlite3", []interface{}{
		models.Task{}, 
		models.TaskSet{},
	})

	if err != nil {
		return nil, err
	}

	persistence.engine = engine
	return &persistence, nil
}

func (p *Persistence) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := p.engine.Select(&tasks, []string{})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (p *Persistence) GetAllTaskSets() ([]models.TaskSet, error) {
	var taskSets []models.TaskSet
	err := p.engine.Select(&taskSets, []string{})
	if err != nil {
		return nil, err
	}
	return taskSets, nil
}

func (p *Persistence) GetTaskByNo(taskNo int64) (models.Task, error) {
	var tasks []models.Task
	err := p.engine.Select(&tasks, []string{fmt.Sprintf("taskNo = %d", taskNo)})
	if err != nil {
		return models.Task{}, err
	}
	return tasks[0], nil
}

func (p *Persistence) GetTaskSetByNo(taskSetNo int64) (models.TaskSet, error) {
	var taskSets []models.TaskSet
	err := p.engine.Select(&taskSets, []string{fmt.Sprintf("taskSetNo = %d", taskSetNo)})
	if err != nil {
		return models.TaskSet{}, err
	}
	return taskSets[0], nil
}

func (p *Persistence) AddTask(task *models.Task) error {
	return p.engine.Insert(task)
}

func (p *Persistence) AddTaskSet(taskSet *models.TaskSet) error {
	return p.engine.Insert(taskSet)
}

func (p *Persistence) UpdateTask(task models.Task) error {
	return p.engine.Update(task)
}

func (p *Persistence) UpdateTaskSet(taskSet models.TaskSet) error {
	return p.engine.Update(taskSet)
}

func (p *Persistence) RemoveTask(task models.Task) error {
	return p.engine.Delete(task)
}

func (p *Persistence) RemoveTaskSet(taskSet models.TaskSet) error {
	return p.engine.Delete(taskSet)
}
