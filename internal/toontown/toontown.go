package toontown

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jaczerob/clerk-gopher/internal/net"
	"github.com/jaczerob/clerk-gopher/internal/sys"
)

var (
	LoginURL     = "https://toontownrewritten.com/api/login"
	LoginHeaders = map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   "clerk-gopher (https://github.com/jaczerob/clerk-gopher)",
	}

	ManifestURL     = "https://cdn.toontownrewritten.com/content/patchmanifest.txt"
	ManifestHeaders = map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "clerk-gopher (https://github.com/jaczerob/clerk-gopher)",
	}

	PatchesURL = "https://download.toontownrewritten.com/patches"

	Platform = sys.GetPlatform()
)

type LoginData struct {
	Success    string `json:"success,omitempty"`
	Gameserver string `json:"gameserver,omitempty"`
	Playcookie string `json:"cookie,omitempty"`
	AppToken   string `json:"appToken,omitempty"`
	AuthToken  string `json:"authToken,omitempty"`
	ETA        string `json:"eta,omitempty"`
	Position   string `json:"position,omitempty"`
	QueueToken string `json:"queueToken,omitempty"`
	Banner     string `json:"banner,omitempty"`
}

type ManifestData struct {
	DL       string                       `json:"dl"`
	Only     []string                     `json:"only"`
	Hash     string                       `json:"hash"`
	CompHash string                       `json:"compHash"`
	Patches  map[string]map[string]string `json:"patches"`
}

func Login(username string, password string) (gameserver string, playcookie string, err error) {
	loginData, err := GetLoginData(username, password)
	if err != nil {
		log.WithField("error", err).Fatal("could not log in")
	}

	gameserver, playcookie = loginData.Gameserver, loginData.Playcookie
	for gameserver == "" && playcookie == "" {
		if loginData.Success == "delayed" {
			log.WithFields(log.Fields{
				"eta":      loginData.ETA,
				"position": loginData.Position,
			}).Println("in queue")

			time.Sleep(5 * time.Second)

			loginData, err := RefreshQueue(loginData.QueueToken)
			if err != nil {
				log.WithField("error", err).Fatal("could not log in")
			}

			gameserver, playcookie = loginData.Gameserver, loginData.Playcookie
		} else if loginData.Success == "false" {
			log.WithField("reason", loginData.Banner).Println("could not log in")
			return
		} else {
			return
		}
	}

	return
}

func GetLoginData(username string, password string) (loginData *LoginData, err error) {
	log.WithField("username", username).Trace("attempting login")

	parameters := map[string]string{
		"username": username,
		"password": password,
		"format":   "json",
	}

	data, err := net.Post(LoginURL, LoginHeaders, parameters)
	if err != nil {
		return nil, fmt.Errorf("could not access login API: %s", err)
	}

	err = json.Unmarshal(data, &loginData)
	if err != nil {
		return nil, fmt.Errorf("error parsing json response: %s\n%s", err, data)
	}

	return loginData, err
}

func RefreshQueue(queueToken string) (loginData *LoginData, err error) {
	log.WithField("queueToken", queueToken).Trace("refreshing queue")

	parameters := map[string]string{
		"queueToken": queueToken,
		"format":     "json",
	}

	data, err := net.Post(LoginURL, LoginHeaders, parameters)
	if err != nil {
		return nil, fmt.Errorf("could not access login API: %s", err)
	}

	err = json.Unmarshal(data, &loginData)
	if err != nil {
		return nil, fmt.Errorf("error parsing json response: %s\n%s", err, data)
	}

	return
}

func GetManifest() (manifest map[string]*ManifestData, err error) {
	log.Trace("getting manifest")

	data, err := net.Get(ManifestURL, ManifestHeaders)
	if err != nil {
		return nil, fmt.Errorf("could not access manifest: %s", err)
	}

	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("error parsing json response: %s\n%s", err, data)
	}

	return
}

func GetHash(filepath string) (string, error) {
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

func PlatformIsIn(only []string) bool {
	for _, platform := range only {
		if platform == Platform {
			return true
		}
	}

	return false
}

func Update(directory string) (err error) {
	log.WithField("dir", directory).Trace("checking for updates")

	manifest, err := GetManifest()
	if err != nil {
		return
	}

	for file, fileManifest := range manifest {
		if !PlatformIsIn(fileManifest.Only) {
			log.WithField("file", file).Trace("ignoring non-OS compliant file")
			continue
		}

		filepath := fmt.Sprintf("%s/%s", directory, file)
		if _, err = os.Stat(filepath); !os.IsNotExist(err) {
			log.WithField("file", file).Trace("exists, hash checking")

			fileHash, err := GetHash(filepath)
			if err != nil {
				return err
			}

			log.WithFields(log.Fields{
				"manifestHash": fileManifest.Hash,
				"fileHash":     fileHash,
			}).Trace("testing hash")

			if fileHash == fileManifest.Hash {
				log.WithField("file", file).Trace("up to date")
				continue
			} else {
				log.WithField("file", file).Trace("out of date")
			}
		} else {
			log.WithField("file", file).Trace("does not exist")
		}

		log.WithField("file", file).Println("updating")

		url := fmt.Sprintf("%s/%s", PatchesURL, fileManifest.DL)
		err = net.Download(filepath, url, ManifestHeaders)
		if err != nil {
			log.WithField("err", err).Error("could not download file")
			return
		}

		log.WithField("file", file).Trace("updated")
	}

	return
}
