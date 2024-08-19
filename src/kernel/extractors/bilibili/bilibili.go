package bilibili

import (
	"duck/kernel/extractors/extractor"
	"duck/kernel/http"
	"duck/kernel/models"
	"duck/kernel/utils"
	"encoding/json"
	"errors"
	"fmt"
)

func init() {
	extractor.Extractors[extractor.Bilibili] = BilibiliExtractor{}
}

const (
	UrlReStr string = `^https://www.bilibili.com/video/BV(.*)/?spm_id_from=(.*)&vd_source=(.*)$`
	BvReStr string = "BV[a-zA-Z0-9]"
	WebInterfaceAPI string = "https://api.bilibili.com/x/web-interface/view"
	PlayerAPI string = "https://api.bilibili.com/x/player/playurl"
)

type BilibiliExtractor struct {}

func (extractor BilibiliExtractor) UrlRe() string {
	return UrlReStr
}

func (extractor BilibiliExtractor) Headers() map[string]string {
	return map[string]string{
		http.Referer: "https://www.bilibili.com",
	}
}

type BilibiliParsedInfo struct {
	listInfo ListInfo
	listCount int
	name string
}

type ListInfo struct {
	bvid string
	aid string
	title string
	pubdate string
	videos []VideoInfo
}

type VideoInfo struct {
	cid string
	title string
	formats []FormatInfo
}

type FormatInfo struct {
	quality string
	format string
	display string
	description string
	urls []string
}

func (extractor BilibiliExtractor) OptionTypes() map[string]string {
	return map[string]string{
		"Name": "Input",
	}
}

func (extractor BilibiliExtractor) Extract(url string) (http.PreflightInfo, interface{}, error) {
	parsedInfo := BilibiliParsedInfo{}
	videoList, err := extractor.extractListInfo(url)
	if err != nil {
		return http.PreflightInfo{}, BilibiliParsedInfo{}, err
	}
	parsedInfo.listInfo = videoList
	parsedInfo.listCount = len(videoList.videos)
	parsedInfo.name = videoList.title
	return http.PreflightInfo{}, parsedInfo, nil
}

func (extractor BilibiliExtractor) extractListInfo(url string) (ListInfo, error) {
	videoList := ListInfo{}

	if res := utils.FindOneSubmatch(BvReStr, url); len(res) != 2 {
		return ListInfo{}, errors.New("no valid bv number was found")
	} else {
		videoList.bvid = utils.FindOneSubmatch(BvReStr, url)[1]
	}

	bytes, err := http.RequestPage(fmt.Sprintf("%s?bvid=%s", WebInterfaceAPI, videoList.bvid))
	if err != nil {
		return ListInfo{}, err
	}
	webInterfaceAPIRes := make(map[string]interface{})
	json.Unmarshal(bytes, &webInterfaceAPIRes)
	if webInterfaceAPIRes["code"] != 0 {
		return ListInfo{}, fmt.Errorf("WebInterfaceAPI returned code was non-zero: %d", webInterfaceAPIRes["code"])
	}
	
	webInterfaceAPIResData := webInterfaceAPIRes["data"].(map[string]interface{})
	videoList.aid = webInterfaceAPIResData["aid"].(string)
	videoList.title = webInterfaceAPIResData["title"].(string)
	videoList.pubdate = webInterfaceAPIResData["pubdate"].(string)

	videoList.videos = make([]VideoInfo, 0)
	for _, page := range webInterfaceAPIResData["pages"].([]map[string]interface{}) {
		vidoeInfo := VideoInfo{}
		vidoeInfo.cid = page["cide"].(string)
		vidoeInfo.title = page["part"].(string)
		formats, err := extractor.parseFormatInfo(videoList.aid, vidoeInfo.cid)
		if err != nil {
			return ListInfo{}, err
		}
		vidoeInfo.formats = formats
		videoList.videos = append(videoList.videos, vidoeInfo)
	}
	return videoList, nil
}

func (extractor BilibiliExtractor) parseFormatInfo(aid string, cid string) ([]FormatInfo, error) {
	bytes, err := http.RequestPage(fmt.Sprintf("%s?avid=%s&cid=%20%s", PlayerAPI, aid, cid))
	if err != nil {
		return nil, err
	}
	playerAPIRes := make(map[string]interface{})
	json.Unmarshal(bytes, &playerAPIRes)
	playerAPIResData := playerAPIRes["data"].(map[string]interface{})

	formatInfos := make([]FormatInfo, 0)
	for _, supportFormat := range playerAPIResData["support_formats"].([]map[string]interface{}) {
		formatInfo := FormatInfo{}
		formatInfo.quality = supportFormat["quality"].(string)
		formatInfo.format = supportFormat["format"].(string)
		formatInfo.description = supportFormat["description"].(string)
		formatInfo.display = supportFormat["display"].(string)
		urls, err := extractor.parseVideoUrls(aid, cid, formatInfo.quality)
		if err != nil {
			return nil, err
		}
		formatInfo.urls = urls
		formatInfos = append(formatInfos, formatInfo)
	}
	return formatInfos, nil
}

func (extractor BilibiliExtractor) parseVideoUrls(aid string, cid string, quality string) ([]string, error) {
	bytes, err := http.RequestPage(fmt.Sprintf("%s?avid=%s&cid=%20%s&qn=%s", PlayerAPI, aid, cid, quality))
	if err != nil {
		return nil, err
	}
	playerAPIRes := make(map[string]interface{})
	json.Unmarshal(bytes, &playerAPIRes)
	playerAPIResData := playerAPIRes["data"].(map[string]interface{})

	res := make([]string, 0)
	for _, item := range playerAPIResData["durl"].([]map[string]interface{}) {
		res = append(res, item["url"].(string))
	}
	return res, nil
}

func (extractor BilibiliExtractor) NewTask(info http.PreflightInfo, options map[string]interface{}) models.Task {
	return models.Task{}
}
