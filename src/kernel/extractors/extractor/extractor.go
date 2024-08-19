package extractor

import (
	"duck/kernel/http"
	"duck/kernel/models"
)

type Extractor interface {
	UrlRe() string
	Headers() map[string]string
	Extract(url string) (http.PreflightInfo, interface{}, error)
	OptionTypes() map[string]string
	NewTask(info http.PreflightInfo, options map[string]interface{}) models.Task
}

var Extractors map[string]Extractor = map[string]Extractor{}

const (
	Http string = "http"
	Bilibili string = "bilibili"
	YouTube string = "YouTube"
)
