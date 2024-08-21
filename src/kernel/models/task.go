package models

import (
	"time"
	"duck/kernel/constance"
)

type Task struct {
	TaskNo int64 `constraints:"PrimaryKey, AutoIncrement"`
	TaskName string	
	Parent int64
	CreateTime time.Time
	UpdateTime time.Time
	FileMimeType string
	FileMimeSubType string
	FileCharset string
	TaskSize int64
	TaskProgress int64
	TaskUrl string
	DownloadUrl string
	TaskStatus constance.TaskStatus
	FileLocation string
	IsRange bool
	Ranges [][]int64
	Extractor string
	AdditionalInfo map[string]string
}
