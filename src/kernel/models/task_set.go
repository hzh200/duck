package models

type TaskSet struct {
	TaskNo int64
	Children []int64
}