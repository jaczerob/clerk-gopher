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

func Get(url string, headers map[string]string) ([]byte, error) {
	log.WithFields(log.Fields{
		"url":     url,
		"headers": headers,
	}).Trace("attempting GET request")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, element := range headers {
		req.Header.Add(key, element)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad http status: %d", resp.StatusCode)
	}

	log.Trace("GET request OK")
	return body, nil
}

func Post(url string, headers map[string]string, parameters map[string]string) ([]byte, error) {
	log.WithFields(log.Fields{
		"url":     url,
		"headers": headers,
	}).Trace("attempting POST request")

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, element := range parameters {
		q.Add(key, element)
	}

	req.URL.RawQuery = q.Encode()

	for key, element := range headers {
		req.Header.Add(key, element)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad http status: %d", resp.StatusCode)
	}

	log.Trace("POST request OK")
	return body, nil
}

func VerifyFilepath(filePath string) (err error) {
	dir := filepath.Dir(filePath)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
	}

	return
}

func Download(filePath string, url string, headers map[string]string) (err error) {
	log.WithFields(log.Fields{
		"filePath": filePath,
		"url":      url,
		"headers":  headers,
	}).Trace("downloading")

	err = VerifyFilepath(filePath)

	if err != nil {
		return
	}

	out, err := os.Create(filePath)
	if err != nil {
		return
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad http status: %d", resp.StatusCode)
	}

	decompressed := bzip2.NewReader(resp.Body)
	_, err = io.Copy(out, decompressed)
	if err != nil {
		return
	}

	return
}
