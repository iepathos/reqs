package reqs

import (
	log "github.com/sirupsen/logrus"
)

func FatalCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
