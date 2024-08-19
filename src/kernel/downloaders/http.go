package downloaders

import (
	"duck/kernel/extractors/extractor"
	"duck/kernel/fs"
	"duck/kernel/http"
	"duck/kernel/models"
	"fmt"
	"io"
	netHttp "net/http"
	"strings"
)

type HttpDownloader struct {
	task *models.Task
	isRunning bool
}

func NewHttpDownloader(task *models.Task) Downloader {	
	downloader := HttpDownloader{}
	downloader.task = task
	downloader.isRunning = false
	return downloader
}

func (downloader HttpDownloader) Start() error {
	req, err := netHttp.NewRequest(netHttp.MethodGet, downloader.task.DownloadUrl, strings.NewReader(""))
	if err != nil {
		return err
	}

	headers, err := http.RequestHeaders(downloader.task.DownloadUrl)
	if err != nil {
		return err
	}

	for header, value := range headers {
		req.Header.Add(header, value)
	}

	if _, found := extractor.Extractors[downloader.task.Extractor]; found {
		for header, value := range extractor.Extractors[downloader.task.Extractor].Headers() {
			req.Header.Add(header, value)
		}
	}

	client := netHttp.Client{}
	// client := &netHttp.Client{
	// 	CheckRedirect: redirectPolicyFunc,
	// }
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	err = fs.CreateIfNotExisted(downloader.task.FileLocation)
    if err != nil {
        return err
    }

	eof := false
	buf := make([]byte, 4096)
	for !eof {
        readN, err := res.Body.Read(buf)
		// EOF can be emitted together with data read.
        if err != nil {
            if err != io.EOF {
				return fmt.Errorf("not reach eof: %v", err)
            }
            eof = true
        }
		writeN, err := fs.WriteFile(downloader.task.FileLocation, downloader.task.TaskProgress, buf[:readN])
		if err != nil {
			return err
		}
		downloader.task.TaskProgress += int64(writeN)
    }

	return nil
}

func (downloader HttpDownloader) Pause() {}

func (downloader HttpDownloader) Done() {}

func (downloader HttpDownloader) Fail() {}
