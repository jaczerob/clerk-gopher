package toontown

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jaczerob/clerk-gopher/internal/net"
)

const loginURL = "https://toontownrewritten.com/api/login"

var loginHeaders = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
	"User-Agent":   "clerk-gopher (https://github.com/jaczerob/clerk-gopher)",
}

type Status struct {
	Open   bool   `json:"open"`
	Banner string `json:"banner,omitempty"`
}

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

func Login(username string, password string) (gameserver string, playcookie string, err error) {
	loginData, err := getLoginData(username, password)
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

			loginData, err := refreshQueue(loginData.QueueToken)
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

func getLoginData(username string, password string) (loginData *LoginData, err error) {
	log.WithField("username", username).Trace("attempting login")

	parameters := map[string]string{
		"username": username,
		"password": password,
		"format":   "json",
	}

	data, err := net.Request("POST", loginURL, loginHeaders, parameters)
	if err != nil {
		return nil, fmt.Errorf("could not access login API: %s", err)
	}

	err = json.Unmarshal(data, &loginData)
	if err != nil {
		return nil, fmt.Errorf("error parsing json response: %s\n%s", err, data)
	}

	return loginData, err
}

func refreshQueue(queueToken string) (loginData *LoginData, err error) {
	log.WithField("queueToken", queueToken).Trace("refreshing queue")

	parameters := map[string]string{
		"queueToken": queueToken,
		"format":     "json",
	}

	data, err := net.Request("POST", loginURL, loginHeaders, parameters)
	if err != nil {
		return nil, fmt.Errorf("could not access login API: %s", err)
	}

	err = json.Unmarshal(data, &loginData)
	if err != nil {
		return nil, fmt.Errorf("error parsing json response: %s\n%s", err, data)
	}

	return
}
