package models

import "time"

type TaskSet struct {
	TaskSetNo int64 `constraints:"PrimaryKey, AutoIncrement"`
	TaskSetName string
	Children []int64
	CreateTime time.Time
	UpdateTime time.Time
	Extractor string
	AdditionalInfo map[string]interface{}
}
