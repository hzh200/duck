package downloaders

import (
	"duck/kernel/extractors/extractor"
	"duck/kernel/http"
	"duck/kernel/log"
	"duck/kernel/models"
	"fmt"
	"io"
	netHttp "net/http"
	"os"
	"strings"
	"sync"
)

type RangeDownloader struct {
	task *models.Task
	fd *os.File
	isRunning bool
}

func NewRangeDownloader(task *models.Task) Downloader {	
	downloader := RangeDownloader{}
	downloader.task = task
	downloader.isRunning = false
	return downloader
}

var mu sync.Mutex

const downloadUnit int64 = 1024 * 1024
const routineCount int = 50

func (downloader RangeDownloader) sliceRange() ([]int64, bool) {
	if len(downloader.task.Ranges) == 0 {
		return nil, false
	}

	r := downloader.task.Ranges[0]
	if r[1] - r[0] <= downloadUnit {
		downloader.task.Ranges = downloader.task.Ranges[1:]
		return r, true
	} else {
		res := []int64{r[0], r[0] + downloadUnit}
		downloader.task.Ranges[0][0] = res[1] + 1
		return res, true
	}
}

func (downloader RangeDownloader) insertRange(r []int64) {
	if len(downloader.task.Ranges) == 0 {
		return
	}

	if r[0] < downloader.task.Ranges[0][0] {
		downloader.task.Ranges = append([][]int64{r}, downloader.task.Ranges...)
		return
	}

	if r[0] > downloader.task.Ranges[len(downloader.task.Ranges) - 1][0] {
		downloader.task.Ranges = append(downloader.task.Ranges, r)
		return
	}

	left, right := 0, len(downloader.task.Ranges) - 1
	for left < right {
		if left == right {
			break
		}
		mid := (right - left) / 2 + left
		if downloader.task.Ranges[mid][0] < r[0] {
			left = mid + 1
		} else {
			right = mid
		}
	}

	downloader.task.Ranges = append(downloader.task.Ranges[:left], append([][]int64{r}, downloader.task.Ranges[left:]...)...)
}

func (downloader RangeDownloader) Start() error {
	return downloader.StartCore()
}

func (downloader *RangeDownloader) StartCore() error {
	downloader.isRunning = true
	var err error
	downloader.fd, err = os.OpenFile(downloader.task.FileLocation, os.O_RDWR | os.O_CREATE, 0666)
    if err != nil {
        return err
    }

	var countMu sync.Mutex
	runningRoutineCount := 0
	
	done := make(chan bool, routineCount)

	go func() {
		for downloader.isRunning && <-done {
			countMu.Lock()
			runningRoutineCount--
			countMu.Unlock()
		}
	}()

	for downloader.isRunning {
		for runningRoutineCount < routineCount {
			slice, existence := downloader.sliceRange()
			if !existence {
				break
			}
			go func() {
				err := downloader.routine(slice[0], slice[1], done)
				if err != nil {
					log.Error(err)
				}
			}()
			countMu.Lock()
			runningRoutineCount++
			countMu.Unlock()
		}

		if runningRoutineCount == 0 && downloader.task.TaskProgress == downloader.task.TaskSize {
			downloader.isRunning = false
		}
	}

	downloader.fd.Close()
	return nil
}

func (downloader RangeDownloader) routine(start, end int64, done chan bool) error {
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
	req.Header.Add(http.Range, fmt.Sprintf("bytes=%d-%d", start, end))

	client := netHttp.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	eof := false
	buf := make([]byte, 4096)
	for !eof {
		readN, err := res.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("not reach eof: %v", err)
            }
            eof = true
		}

		// writeN, err := fs.WriteFile(downloader.task.FileLocation, start, buf[:readN])
		writeN, err := downloader.fd.WriteAt(buf[:readN], start)
		if err != nil {
			return err
		}
		if writeN != readN {
			return fmt.Errorf("readN: %d, writeN: %d", readN, writeN)
		}

		start += int64(writeN)
		downloader.updateProgress(int64(writeN))
	}

	done <- true
	return nil
}

func (downloader RangeDownloader) updateProgress(progress int64) {
	mu.Lock()
	downloader.task.TaskProgress += int64(progress)
	mu.Unlock()
}

func (downloader RangeDownloader) Pause() {
	
}

func (downloader RangeDownloader) Done() {
	
}

func (downloader RangeDownloader) Fail() {
	
}
