package http

import (
	"duck/kernel/utils"
	"errors"
	"net/http"
	"strings"
)

const (
	UrlReStr = `^(https?):\/\/[-A-Za-z0-9+&@#\/%?=~_|!:,.;]+[-A-Za-z0-9+&@#\/%=~_|]`
	RedirectLimit = 4
)

func combineRelativePath(basePath string, relativePath string) string {
	relativePath = strings.TrimPrefix(relativePath, "/")
	if !strings.Contains(basePath, "/") {
		if strings.Contains(basePath, "?") {
			basePath = strings.Split(basePath, "?")[0]
		}
		return basePath + "/" + relativePath
	}
	if strings.Contains(relativePath, "/") {
		overlap := strings.Split(relativePath, "/")[0]
		return basePath[:strings.Index(basePath, overlap)] + relativePath
	} else {
		return basePath[:strings.LastIndex(basePath, "/")] + "/" + relativePath
	}
}

func checkIsRedirected(res *http.Response) bool {
	return res.StatusCode == 302 || res.StatusCode == 307 || res.StatusCode == 308
}

func parseRedirectUrl(req *http.Request, res *http.Response) (string, error) {
	if headers, found := res.Header[Location]; found {
		return headers[0], nil
	}
	data, err := RequestPage(req.URL.String())
	if err != nil {
		return "", err
	}
	if match := utils.FindOneSubmatch(UrlReStr, string(data)); len(match) != 2 {
		return "", nil
	} else {
		return combineRelativePath(req.URL.String(), match[1]), nil
	}
}

func handleRedirect(method string, url string, additionalHeaders map[string]string) (string, *http.Response, error) {
	clt := http.Client{}
	var req *http.Request
	var res *http.Response
	var err error
	redirectCount := 0
	for  {
		if redirectCount >= RedirectLimit {
			return "", res, errors.New("reached the maximum redirection count")
		}

		req, err = http.NewRequest(method, url, strings.NewReader(""))
		if err != nil {
			return "", res, err
		}

		headers, err := RequestHeaders(url)
		if err != nil {
			return "", res, err
		}
	
		for header, value := range headers {
			req.Header.Add(header, value)
		}

		for header, value := range additionalHeaders {
			req.Header.Add(header, value)
		}

		res, err = clt.Do(req)
		if err != nil {
			return "", res, err
		}
		
		if !checkIsRedirected(res) {
			break
		}

		url, err = parseRedirectUrl(req, res)
		if err != nil {
			return "", res, err
		}
		
		redirectCount++
	}
	
	return url, res, nil
}
