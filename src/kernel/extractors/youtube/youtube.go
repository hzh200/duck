package extractors

import (
	"duck/kernel/extractors/extractor"
	"duck/kernel/http"
	"duck/kernel/models"
	"duck/kernel/utils"
	"encoding/json"
	"errors"
	netUrl "net/url"
	"strings"
)

func init() {
	extractor.Extractors[extractor.YouTube] = YouTubeExtractor{}
}

const (
	Host string = "https://www.youtube.com"
	UrlReStr string = `^https://www.youtube.com/watch?v=(.*)$`
	TitleReStr string = `<title>(.*) - YouTube<\/title>`
	MutiplexedFormatsReStr string = `\"formats\".+?]`
	AdaptiveFormatsReStr string = `\"adaptiveFormats\".+?]`
	MimetypeReStr string = `\(video\|audio\)\/(.+); codecs=\"(.+)\"`
	Html5PlayerReStr string = `\"(\/s\/player.+?)\"`
)

type YouTubeExtractor struct {}

func (extractor YouTubeExtractor) UrlRe() string {
	return UrlReStr
}

func (extractor YouTubeExtractor) Headers() map[string]string {
	return map[string]string{
		http.Referer: Host,
	}
}

type YouTubeParsedInfo struct {
	name string
	multiplexedFormats []FormatInfo
	videoFormats []FormatInfo
	audioFormats []FormatInfo
}

func (extractor YouTubeExtractor) OptionTypes() map[string]string {
	return map[string]string{
		"Name": "Input",
	}
}

type FormatInfo struct {
    itag string
    url string
    quality string
    mimeType string
    mimeSubType string
    codecs string
    publishedTimestamp string
}

var CodecsMap = map[string]string{
	"vp9": "vp9",
	"av01": "av1",
	"avc": "avc",
	"opus": "opus",
	"mp4a": "mp4a",
}

func (extractor YouTubeExtractor) Extract(url string) (http.PreflightInfo, interface{}, error) {
	parsedUrlInfo, err := netUrl.Parse(url)
	if err != nil {
		return http.PreflightInfo{}, YouTubeParsedInfo{}, err
	}

	params := parsedUrlInfo.Query()
	params.Add("bpctr", "9999999999")
	params.Add("has_verified", "1")
	parsedUrlInfo.RawQuery = params.Encode()
	bytes, err := http.RequestPage(parsedUrlInfo.String())
	if err != nil {
		return http.PreflightInfo{}, YouTubeParsedInfo{}, err
	}
	pageData := string(bytes)

	var playerUrl string
	if res := utils.FindOneSubmatch(Html5PlayerReStr, pageData); len(res) != 2 {
		return http.PreflightInfo{}, YouTubeParsedInfo{}, errors.New("no valid player url was found")
	} else {
		playerUrl = res[1]
	}
	
	bytes, err = http.RequestPage(playerUrl)
	if err != nil {
		return http.PreflightInfo{}, YouTubeParsedInfo{}, err
	}
	// html5Player := string(bytes)

	parsedInfo := YouTubeParsedInfo{}

	if res := utils.FindOneSubmatch(TitleReStr, pageData); len(res) != 2 {
		return http.PreflightInfo{}, YouTubeParsedInfo{}, errors.New("no valid player url was found")
	} else {
		parsedInfo.name = res[1]
	}


	parsedInfo.multiplexedFormats = make([]FormatInfo, 0)
	parsedInfo.videoFormats = make([]FormatInfo, 0)
	parsedInfo.audioFormats = make([]FormatInfo, 0)
	
	var multiplexedFormatsJSONData string
	var adaptiveFormatsJSONData string
	for _, line := range strings.Split(pageData, "\n") {
		if strings.Contains(line, "streamingData") {
			if res := utils.FindAllSubmatch(MutiplexedFormatsReStr, line); len(res) != 2 {
				return http.PreflightInfo{}, YouTubeParsedInfo{}, errors.New("no valid multiplexed format was found")
			} else {
				multiplexedFormatsJSONData = "{" + res[len(res) - 1][0] + "}"
			}

			if res := utils.FindAllSubmatch(AdaptiveFormatsReStr, line); len(res) != 2 {
				return http.PreflightInfo{}, YouTubeParsedInfo{}, errors.New("no valid adaptive format was found")
			} else {
				adaptiveFormatsJSONData = "{" + res[0][0] + "}"
			}
		}
	}

	multiplexedFormatsRes := make(map[string]interface{})
	adaptiveFormatsRes := make(map[string]interface{})
	json.Unmarshal([]byte(multiplexedFormatsJSONData), &multiplexedFormatsRes)
	json.Unmarshal([]byte(adaptiveFormatsJSONData), &adaptiveFormatsRes)
	for _, multiplexedFormat := range multiplexedFormatsRes["formats"].([]map[string]interface{}) {
		if multiplexedFormat["url"] == nil && multiplexedFormat["signatureCipher"] == nil {
			continue
		}

		formatInfo := FormatInfo{}
		if multiplexedFormat["url"] != nil {
			// decipherN(html5Player, multiplexedFormat["url"])
		} else {
			// decipherSignature(html5player, multiplexedFormat["signatureCipher"])
		}

		if res := utils.FindOneSubmatch(MimetypeReStr, multiplexedFormat["mimeType"].(string)); len(res) != 3 {
			return http.PreflightInfo{}, YouTubeParsedInfo{}, errors.New("no valid mime type was found")
		} else {
			formatInfo.mimeType = res[1]
			formatInfo.mimeSubType = res[2]
			formatInfo.codecs = res[2]
		}

		formatInfo.quality = multiplexedFormat["qualityLabel"].(string)
		formatInfo.publishedTimestamp = multiplexedFormat["lastModified"].(string)
		parsedInfo.multiplexedFormats = append(parsedInfo.multiplexedFormats, formatInfo)
	}

	for _, adaptiveFormat := range adaptiveFormatsRes["adaptiveFormats"].([]map[string]interface{}) {
		if adaptiveFormat["url"] == nil && adaptiveFormat["signatureCipher"] == nil {
			continue
		}

		formatInfo := FormatInfo{}

		formatInfo.itag = adaptiveFormat["itag"].(string)

		if adaptiveFormat["url"] != nil {
			// decipherN(html5Player, adaptiveFormat["url"])
		} else {
			// decipherSignature(html5player, adaptiveFormat["signatureCipher"])
		}

		if res := utils.FindOneSubmatch(MimetypeReStr, adaptiveFormat["mimeType"].(string)); len(res) != 3 {
			return http.PreflightInfo{}, YouTubeParsedInfo{}, errors.New("no valid mime type was found")
		} else {
			formatInfo.mimeType = res[1]
			formatInfo.mimeSubType = res[2]
			for prefix, codecs := range CodecsMap {
				if strings.HasPrefix(res[3], prefix) {
					formatInfo.codecs = codecs
					break
				}
			}
		}

		formatInfo.publishedTimestamp = adaptiveFormat["lastModified"].(string)
		
		if formatInfo.mimeType == "video" {
			formatInfo.quality =  adaptiveFormat["qualityLabel"].(string)
			lastFormatInfo := parsedInfo.videoFormats[len(parsedInfo.videoFormats) - 1]
			if len(parsedInfo.videoFormats) != 0 && 
				!(lastFormatInfo.codecs == formatInfo.codecs && lastFormatInfo.quality == formatInfo.quality && lastFormatInfo.mimeSubType == formatInfo.mimeSubType) {
					parsedInfo.videoFormats = parsedInfo.videoFormats[:len(parsedInfo.videoFormats) - 1]
			}
			parsedInfo.videoFormats = append(parsedInfo.videoFormats, formatInfo)
		} else if formatInfo.mimeType == "audio" {
			formatInfo.quality =  adaptiveFormat["qualityLabel"].(string)
			lastFormatInfo := parsedInfo.audioFormats[len(parsedInfo.audioFormats) - 1]
			if len(parsedInfo.audioFormats) != 0 && 
				!(lastFormatInfo.codecs == formatInfo.codecs && lastFormatInfo.quality == formatInfo.quality && lastFormatInfo.mimeSubType == formatInfo.mimeSubType) {
					parsedInfo.audioFormats = parsedInfo.audioFormats[:len(parsedInfo.audioFormats) - 1]
			}
			parsedInfo.audioFormats = append(parsedInfo.audioFormats, formatInfo)
		}		
	}
	return http.PreflightInfo{}, parsedInfo, nil
}

func (extractor YouTubeExtractor) NewTask(info http.PreflightInfo, options map[string]interface{}) models.Task {
	return models.Task{}
}
