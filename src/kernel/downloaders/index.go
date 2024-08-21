package downloaders

import (
	"duck/kernel/models"
)

type Downloader interface {
	Start() error
	Pause()
	Done()
	Fail()
}

func Download(task *models.Task) error {
	var downloader Downloader
	if !task.IsRange {
		downloader = NewHttpDownloader(task)
	} else {
		downloader = NewRangeDownloader(task)
	}
	return downloader.Start()
}
