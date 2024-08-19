package extractors

import (
	_ "duck/kernel/extractors/bilibili"
	"duck/kernel/extractors/extractor"
	_ "duck/kernel/extractors/http"
	_ "duck/kernel/extractors/youtube"
	"duck/kernel/http"
	"duck/kernel/models"
	"duck/kernel/utils"
	"errors"
)

var infoCache map[string]http.PreflightInfo = make(map[string]http.PreflightInfo)

func Extract(url string, extractorName string) (interface{}, map[string]string, error) {
	// for _, extractor := range Extractors {
	// 	if utils.Match(extractor.UrlRe(), url) {
	// 		return extractor.Extract(url)
	// 	}
	// }
	extractor, found := extractor.Extractors[extractorName]
	if !found {
		return nil, nil, errors.New("extractor doesn't exist")
	}
	if !utils.Match(extractor.UrlRe(), url) {
		return nil, nil, errors.New("url not match")
	}
	info, options, err := extractor.Extract(url)
	optionTypes := extractor.OptionTypes()
	infoCache[extractorName + "|" + url] = info
	return options, optionTypes, err
}

func NewTask(url string, extractorName string, options map[string]interface{}) (models.Task, error) {
	extractor, found := extractor.Extractors[extractorName]
	if !found {
		return models.Task{}, errors.New("")
	}
	info, found := infoCache[extractorName + "|" + url]
	if !found {
		return models.Task{}, errors.New("")
	}
	return extractor.NewTask(info, options), nil
}
