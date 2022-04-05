package net

import (
	"compress/bzip2"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func makeRequest(method string, url string, headers map[string]string, parameters map[string]string) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		return
	}

	for key, element := range headers {
		req.Header.Add(key, element)
	}

	if parameters != nil {
		q := req.URL.Query()
		for key, element := range parameters {
			q.Add(key, element)
		}
		req.URL.RawQuery = q.Encode()
	}

	return
}

func Request(method string, url string, headers map[string]string, parameters map[string]string) (body []byte, err error) {
	log.WithFields(log.Fields{
		"method":     method,
		"url":        url,
		"headers":    headers,
		"parameters": parameters,
	}).Trace("attempting request")

	req, err := makeRequest(method, url, headers, parameters)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad http status: %d", resp.StatusCode)
	}

	log.Trace("request OK")
	return
}

func getFile(filePath string) (file *os.File, err error) {
	dir := filepath.Dir(filePath)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		errDir := os.MkdirAll(dir, 0700)
		if errDir != nil {
			return nil, errDir
		}
	}

	return os.Create(filePath)
}

func getURLBody(url string, headers map[string]string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad http status: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func DownloadBZ2(filePath string, url string, headers map[string]string) (err error) {
	log.WithFields(log.Fields{
		"filePath": filePath,
		"url":      url,
		"headers":  headers,
	}).Trace("downloading")

	out, err := getFile(filePath)
	if err != nil {
		return
	}

	defer out.Close()

	body, err := getURLBody(url, headers)
	if err != nil {
		return
	}

	decompressed := bzip2.NewReader(body)
	_, err = io.Copy(out, decompressed)
	if err != nil {
		return
	}

	return
}
