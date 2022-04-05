package main

import (
	"github.com/jaczerob/clerk-gopher/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Panic(err)
	}
}
