package update

import (
	"compress/bzip2"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/jaczerob/clerk-gopher/internal/static"
	"github.com/jaczerob/clerk-gopher/internal/sys"
	log "github.com/sirupsen/logrus"
)

type UpdateClient struct {
	Directory string

	http    *http.Client
	headers map[string]string
	baseURL *url.URL
}

func NewUpdateClient() (c *UpdateClient, err error) {
	baseURL, err := url.Parse(static.UpdateDownloadURL)
	if err != nil {
		return
	}

	dir, err := sys.GetDirectory()
	if err != nil {
		return
	}

	c = &UpdateClient{
		http:      &http.Client{},
		headers:   static.Headers,
		baseURL:   baseURL,
		Directory: dir,
	}

	return
}

func (c *UpdateClient) Update() (err error) {
	updateFiles, err := c.getUpdateFiles()
	if err != nil {
		return
	}

	for _, file := range updateFiles {
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

	out, err := getFile(file.Path)
	if err != nil {
		return
	}

	defer out.Close()

	resp, err := c.http.Get(file.URL)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad http status: %d", resp.StatusCode)
	}

	decompressed := bzip2.NewReader(resp.Body)
	_, err = io.Copy(out, decompressed)
	return
}

func (c *UpdateClient) getUpdateFiles() (updateFiles []*UpdateFile, err error) {
	log.WithField("dir", c.Directory).Trace("checking for updates")

	manifest, err := c.getManifest()
	if err != nil {
		return
	}

	for file, fileManifest := range manifest {
		if !platformIsIn(fileManifest.Only) {
			log.WithField("file", file).Trace("ignoring non-OS compliant file")
			continue
		}

		filepath := path.Join(c.Directory, file)
		if _, statErr := os.Stat(filepath); !os.IsNotExist(statErr) {
			log.WithField("file", file).Trace("exists, hash checking")

			fileHash, fileErr := getHash(filepath)
			if fileErr != nil {
				return nil, fileErr
			}

			log.WithFields(log.Fields{
				"manifestHash": fileManifest.Hash,
				"fileHash":     fileHash,
			}).Trace("testing hash")

			if fileHash == fileManifest.Hash {
				log.WithField("file", file).Trace("up to date")
				continue
			} else {
				log.WithField("file", file).Trace("out of date, queueing for update")
			}
		} else {
			log.WithField("file", file).Trace("does not exist, queueing for download")
		}

		u := *c.baseURL
		u.Path = path.Join(u.Path, fileManifest.DL)

		updateFiles = append(updateFiles, &UpdateFile{
			Name: file,
			Path: filepath,
			URL:  u.String(),
		})
	}

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

func getHash(filepath string) (string, error) {
	log.WithField("file", filepath).Trace("getting file hash")

	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func platformIsIn(only []string) bool {
	curPlatform := sys.GetPlatform()

	for _, platform := range only {
		if platform == curPlatform {
			return true
		}
	}

	return false
}
