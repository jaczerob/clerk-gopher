package toontown

import (
	"compress/bzip2"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/jaczerob/clerk-gopher/internal/net"
	"github.com/jaczerob/clerk-gopher/internal/sys"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

const (
	manifestURL = "https://cdn.toontownrewritten.com/content/patchmanifest.txt"
	patchesURL  = "https://download.toontownrewritten.com/patches"
)

var manifestHeaders = map[string]string{
	"Content-Type": "application/json",
	"User-Agent":   "clerk-gopher (https://github.com/jaczerob/clerk-gopher)",
}

type ManifestData struct {
	DL       string                       `json:"dl"`
	Only     []string                     `json:"only"`
	Hash     string                       `json:"hash"`
	CompHash string                       `json:"compHash"`
	Patches  map[string]map[string]string `json:"patches"`
}

func getManifest() (manifest map[string]*ManifestData, err error) {
	log.Trace("getting manifest")

	data, err := net.Request("GET", manifestURL, manifestHeaders, nil)
	if err != nil {
		return nil, fmt.Errorf("could not access manifest: %s", err)
	}

	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("error parsing json response: %s\n%s", err, data)
	}

	return
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

type UpdateFile struct {
	Name string
	Path string
	URL  string
}

func getUpdateFiles(directory string) (updateFiles []*UpdateFile, err error) {
	log.WithField("dir", directory).Trace("checking for updates")

	manifest, err := getManifest()
	if err != nil {
		return
	}

	for file, fileManifest := range manifest {
		if !platformIsIn(fileManifest.Only) {
			log.WithField("file", file).Trace("ignoring non-OS compliant file")
			continue
		}

		filepath := path.Join(directory, file)
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

		updateFiles = append(updateFiles, &UpdateFile{
			Name: file,
			Path: filepath,
			URL:  fmt.Sprintf("%s/%s", patchesURL, fileManifest.DL),
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

func downloadFile(f *UpdateFile) (err error) {
	req, err := http.NewRequest("GET", f.URL, nil)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	out, err := getFile(f.Path)
	if err != nil {
		return
	}

	defer out.Close()

	decompressReader := bzip2.NewReader(resp.Body)
	bar := progressbar.DefaultBytes(-1, fmt.Sprintf("downloading %s", f.Name))

	_, err = io.Copy(io.MultiWriter(out, bar), decompressReader)
	if err != nil {
		return
	}

	fmt.Println()
	return
}

func Update(directory string) (err error) {
	updateFiles, err := getUpdateFiles(directory)
	if err != nil {
		return
	}

	for _, file := range updateFiles {
		err := downloadFile(file)
		if err != nil {
			return err
		}
	}

	return
}
