package login

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/jaczerob/clerk-gopher/internal/static"
	log "github.com/sirupsen/logrus"
)

type LoginClient struct {
	http    *http.Client
	headers map[string]string
	baseURL *url.URL
}

func NewLoginClient() *LoginClient {
	baseURL, _ := url.Parse(static.LoginEndpoint)
	return &LoginClient{
		http:    &http.Client{Timeout: 5 * time.Second},
		headers: static.Headers,
		baseURL: baseURL,
	}
}

func (c *LoginClient) Login(username string, password string) (gameserver string, playcookie string, err error) {
	loginData, err := c.getLoginData(username, password)
	if err != nil {
		log.WithField("error", err).Fatal("could not log in")
	}

	gameserver, playcookie = loginData.Gameserver, loginData.Playcookie
	for gameserver == "" && playcookie == "" {
		if loginData.Success == static.LoginSuccessDelayed {
			log.WithFields(log.Fields{
				"eta":      loginData.ETA,
				"position": loginData.Position,
			}).Println("in queue")

			time.Sleep(5 * time.Second)

			loginData, err := c.refreshQueue(loginData.QueueToken)
			if err != nil {
				log.WithField("error", err).Fatal("could not log in")
			}

			gameserver, playcookie = loginData.Gameserver, loginData.Playcookie
		} else if loginData.Success == static.LoginSuccessFailed {
			log.WithField("reason", loginData.Banner).Println("could not log in")
			return
		} else {
			return
		}
	}

	return
}

func (c *LoginClient) getLoginData(username string, password string) (loginData *LoginData, err error) {
	log.WithField("username", username).Trace("attempting login")

	user := url.Values{
		"username": {username},
		"password": {password},
	}

	data, err := c.post(user)
	if err != nil {
		return nil, fmt.Errorf("could not access login API: %s", err)
	}

	err = json.Unmarshal(data, &loginData)
	if err != nil {
		return nil, fmt.Errorf("error parsing json response: %s\n%s", err, data)
	}

	return loginData, err
}

func (c *LoginClient) refreshQueue(queueToken string) (loginData *LoginData, err error) {
	log.WithField("queueToken", queueToken).Trace("refreshing queue")

	queuedUser := url.Values{
		"queueToken": {queueToken},
	}

	data, err := c.post(queuedUser)
	if err != nil {
		return nil, fmt.Errorf("could not access login API: %s", err)
	}

	err = json.Unmarshal(data, &loginData)
	if err != nil {
		return nil, fmt.Errorf("error parsing json response: %s\n%s", err, data)
	}

	return
}

func (c *LoginClient) post(values url.Values) (body []byte, err error) {
	resp, err := c.http.PostForm(c.baseURL.String(), values)
	if err != nil {
		return
	}

	log.WithFields(log.Fields{
		"method": resp.Request.Method,
		"url":    resp.Request.URL,
	}).Trace("request made")

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
