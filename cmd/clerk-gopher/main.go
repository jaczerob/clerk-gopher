package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jaczerob/clerk-gopher/internal/sys"
	"github.com/jaczerob/clerk-gopher/internal/toontown"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	Username string
	Password string
	Verbose  bool
)

func init() {
	flag.StringVar(&Username, "username", "", "TTR username")
	flag.StringVar(&Password, "password", "", "TTR password")
	flag.BoolVar(&Verbose, "v", false, "Enable verbose logging")
	flag.Parse()

	if Verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&prefixed.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
		DisableSorting:  true,
	})
}

func main() {
	directory, err := sys.GetDirectory()
	if err != nil {
		log.WithField("error", err).Fatal("could not initialize toontown rewritten directory")
	}

	executable := fmt.Sprintf("%s/%s", directory, sys.GetExecutable())

	err = toontown.Update(directory)
	if err != nil {
		log.WithField("error", err).Fatal("could not update")
	}

	loginData, err := toontown.Login(Username, Password)
	if err != nil {
		log.WithField("error", err).Fatal("could not log in")
	}

	gameserver, cookie := loginData.Gameserver, loginData.Cookie
	for gameserver == "" && cookie == "" {
		if loginData.Success == "delayed" {
			log.WithFields(log.Fields{
				"eta":      loginData.ETA,
				"position": loginData.Position,
			}).Println("in queue")

			time.Sleep(5 * time.Second)

			loginData, err := toontown.RefreshQueue(loginData.QueueToken)
			if err != nil {
				log.WithField("error", err).Fatal("could not log in")
			}

			gameserver, cookie = loginData.Gameserver, loginData.Cookie
		} else if loginData.Success == "false" {
			log.WithField("reason", loginData.Banner).Println("could not log in")
			return
		} else {
			return
		}
	}

	log.Println("entering toontown, have fun!")
	sys.RunExecutable(executable, gameserver, cookie)
}
