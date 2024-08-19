package http

import (
	"duck/kernel/utils"
	"fmt"
	"net/http"
	netUrl "net/url"
	"strconv"
	"strings"
	"time"
)

type PreflightInfo struct {
	Url string
	Name string
	Size int64
	MimeType string
    MimeSubType string
	Charset string
	PubDate time.Time
	DownloadUrl string
	IsRange bool
}

func Preflight(url string) (PreflightInfo, error) {
	redirectedUrl, res, err := handleRedirect(http.MethodGet, url, map[string]string{Range: "bytes=0-0"})
    if err != nil {
        return PreflightInfo{}, err
    }
	
	if !(res.StatusCode == 200 || res.StatusCode == 206) {
		return PreflightInfo{}, fmt.Errorf("preflight status code %d", res.StatusCode)
	}
    res.Body.Close()

	preflightInfo := PreflightInfo{}
	preflightInfo.Url = url
	preflightInfo.DownloadUrl = redirectedUrl
	if headers, found := res.Header[ContentType]; found {
		match := utils.FindOneSubmatch(MimetypeReStr, headers[0])
		if len(match) >= 3 {
			preflightInfo.MimeType = match[1]
			preflightInfo.MimeSubType = match[2]
		}
		if len(match) > 3 {
			preflightInfo.Charset = match[3]
		}
	}

	if headers, found := res.Header[ContentDisposition]; found {
		match := utils.FindOneSubmatch(`filename=\"(.+)\"`, headers[0])
		if len(match) == 2 {
			preflightInfo.Name = match[1]
		}
	}

	if preflightInfo.Name == "" {
		if strings.Contains(url, "//") {
			url = strings.Split(url, "//")[1]
		}
		if strings.Contains(url, "/") && !strings.HasSuffix(url, "/") {
			parsedUrl, err := netUrl.Parse(url)
			if err != nil {
				return PreflightInfo{}, err
			}

			urlParts := strings.Split(parsedUrl.Path, "/")
			preflightInfo.Name = urlParts[len(urlParts) - 1]
			if !strings.Contains(preflightInfo.Name, ".") {
				preflightInfo.Name = preflightInfo.Name + "." + preflightInfo.MimeSubType
			}
		} else {
			preflightInfo.Name = preflightInfo.MimeType + "." + preflightInfo.MimeSubType
		}
	}

	if res.StatusCode == 200 {
		if headers, found := res.Header[ContentLength]; found {
			preflightInfo.Size, err = strconv.ParseInt(headers[0], 0, 64)
			if err != nil {
				return PreflightInfo{}, err
			}
		}
	} else {
		if headers, found := res.Header[ContentRange]; found {
			parts := strings.Split(headers[0], "/")
			preflightInfo.Size, err = strconv.ParseInt(parts[1], 0, 64)
			if err != nil {
				return PreflightInfo{}, err
			}
		}
	}

	if headers, found := res.Header[LastModified]; found {
		match := utils.FindOneSubmatch(LastModifiedReStr, headers[0])
		if len(match) == 10 {
			year, err := strconv.ParseInt(match[4], 0, 64)
			if err != nil {
				return PreflightInfo{}, err
			}
			var month time.Month
			switch match[3] {
			case "Jan":
				month = time.January
			case "Feb":
				month = time.February
			case "Mar":
				month = time.March
			case "Apr":
				month = time.April
			case "May":
				month = time.May
			case "Jun":
				month = time.June
			case "Jul":
				month = time.July
			case "Aug":
				month = time.August
			case "Sep":
				month = time.September
			case "Oct":
				month = time.October
			case "Nov":
				month = time.November
			case "Dec":
				month = time.December
			}
			day, err := strconv.ParseInt(match[2], 0, 64)
			if err != nil {
				return PreflightInfo{}, err
			}
			hour, err := strconv.ParseInt(match[5], 0, 64)
			if err != nil {
				return PreflightInfo{}, err
			}
			min, err := strconv.ParseInt(match[6], 0, 64)
			if err != nil {
				return PreflightInfo{}, err
			}
			sec, err := strconv.ParseInt(match[7], 0, 64)
			if err != nil {
				return PreflightInfo{}, err
			}
			preflightInfo.PubDate = time.Date(int(year), month, int(day), int(hour), int(min), int(sec), 0, time.UTC)
		}
	} else {
		preflightInfo.PubDate = time.Now()
	}

	preflightInfo.IsRange = (res.StatusCode == 206)

	return preflightInfo, nil
}
