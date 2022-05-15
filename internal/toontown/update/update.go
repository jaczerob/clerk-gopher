package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/jaczerob/clerk-gopher/internal/static"
	"github.com/jaczerob/clerk-gopher/internal/util"
	log "github.com/sirupsen/logrus"
)

func NewUpdateClient(executable *util.Executable) (c *UpdateClient, err error) {
	baseURL, err := url.Parse(static.UpdateDownloadURL)
	if err != nil {
		return
	}

	c = &UpdateClient{
		executable: executable,
		http:       &http.Client{},
		headers:    static.Headers,
		baseURL:    baseURL,
	}

	return
}

func (c *UpdateClient) Update() (err error) {
	updateFiles, err := c.getUpdateFiles()
	if err != nil {
		return
	}

	for _, file := range updateFiles {
		log.Info("updating", file.Name)

		err := c.downloadBZ2(file)
		if err != nil {
			return err
		}
	}

	return
}

func (c *UpdateClient) getManifest() (manifest map[string]*ManifestData, err error) {
	log.Trace("getting manifest")

	data, err := c.get(static.UpdateManifestURL)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("error parsing json response: %s\n%s", err, data)
	}

	return
}

func (c *UpdateClient) get(url string) (body []byte, err error) {
	resp, err := c.http.Get(url)
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

func (c *UpdateClient) downloadBZ2(file *UpdateFile) (err error) {
	log.WithFields(log.Fields{
		"name":     file.Name,
		"filePath": file.Path,
		"url":      file.URL,
	}).Trace("downloading")

	resp, err := c.http.Get(file.URL)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad http status: %d", resp.StatusCode)
	}

	file.Read(resp.Body)
	return
}

func (c *UpdateClient) getUpdateFiles() (updateFiles []*UpdateFile, err error) {
	log.Trace("checking for updates")

	manifest, err := c.getManifest()
	if err != nil {
		return
	}

	for file, fileManifest := range manifest {
		if !c.platformIsIn(fileManifest.Only) {
			log.WithField("file", file).Trace("ignoring non-OS compliant file")
			continue
		}

		u := *c.baseURL
		u.Path = path.Join(u.Path, fileManifest.DL)

		filepath := path.Join(c.executable.Directory, file)
		updateFile, err := NewUpdateFile(file, filepath, u.String())
		if err != nil {
			return nil, err
		}

		if updateFile.Exists {
			log.WithFields(log.Fields{
				"file":         file,
				"manifestHash": fileManifest.Hash,
				"fileHash":     updateFile.Hash,
			}).Trace("exists, testing hash")

			if updateFile.Hash == fileManifest.Hash {
				log.WithField("file", file).Trace("up to date")
				continue
			} else {
				log.WithField("file", file).Trace("out of date, queueing for update")
			}
		} else {
			log.WithField("file", file).Trace("does not exist, queueing for download")
		}

		updateFiles = append(updateFiles, updateFile)
	}

	return
}

func (c *UpdateClient) platformIsIn(only []string) bool {
	for _, platform := range only {
		if platform == c.executable.Platform {
			return true
		}
	}

	return false
}
