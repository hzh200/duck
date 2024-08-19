package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func RequestPage(url string) ([]byte, error) {
    req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
    if err != nil {
        return []byte{}, err
    }

    headers, err := RequestHeaders(url)
    if err != nil {
        return []byte{}, err
    }

    for header, value := range headers {
        req.Header.Add(header, value)
    }

    resp, err := http.Get(url)
    if err != nil {
        return []byte{}, err
    }
    defer resp.Body.Close()

    buf := make([]byte, 4096)
	bytes := make([]byte, 0)
	for {
        n, err := resp.Body.Read(buf)
        if err != nil {
            if err != io.EOF {
				return []byte{}, fmt.Errorf("not reach eof: %v", err)
            }
            break
        }
		
		bytes = append(bytes, buf[:n]...)
    }

	return bytes, nil
}
