package helper

import (
	"log"

	"github.com/stretchr/testify/suite"
)

func ErrorLog(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ErrorIts(err error) {
	var its suite.Suite
	if err != nil {
		its.FailNow(err.Error())
	}
}
