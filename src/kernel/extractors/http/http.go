package http

import (
	"duck/kernel/constance"
	"duck/kernel/extractors/extractor"
	"duck/kernel/http"
	"duck/kernel/models"
	"path"
	"time"
)

func init() {
	extractor.Extractors[extractor.Http] = HttpExtractor{}
}

const (
	UrlReStr string = http.UrlReStr
)

type HttpExtractor struct {}

func (extractor HttpExtractor) UrlRe() string {
	return UrlReStr
}

func (extractor HttpExtractor) Headers() map[string]string {
	return map[string]string{}
}

type HttpParsedInfo struct {
	Name string
}

func (extractor HttpExtractor) OptionTypes() map[string]string {
	return map[string]string{
		"Name": "Input",
	}
}

func (extractor HttpExtractor) Extract(url string) (http.PreflightInfo, interface{}, error) {
	preflightInfo, err := http.Preflight(url)
	if err != nil {
		return http.PreflightInfo{}, HttpParsedInfo{}, err
	}

	httpParsedInfo := HttpParsedInfo{}
	httpParsedInfo.Name = preflightInfo.Name
	return preflightInfo, httpParsedInfo, nil
}

func (extractor HttpExtractor) NewTask(info http.PreflightInfo, options map[string]interface{}) models.Task {
	task := models.Task{}
	task.TaskName = options["Name"].(string)
	task.CreateTime = info.PubDate
	task.UpdateTime = time.Now()
	task.FileMimeType = info.MimeType
	task.FileCharset = info.Charset
	task.TaskSize = info.Size
	task.TaskProgress = 0
	task.TaskUrl = info.Url
	task.DownloadUrl = info.DownloadUrl
	task.TaskStatus = constance.Waiting
	task.FileLocation = path.Join(options["Location"].(string), task.TaskName)
	task.IsRange = info.IsRange
	task.Ranges = [][]int64{{0, task.TaskSize - 1}}
	task.Extractor = options["Extractor"].(string)
	task.AdditionalInfo = make(map[string]string)
	return task
}
