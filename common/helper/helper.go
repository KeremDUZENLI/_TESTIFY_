package helper

import (
	"log"

	"github.com/stretchr/testify/suite"
)

func ErrorPrint(err error) {
	if err != nil {
		println(err.Error())
	}
}

func ErrorLog(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ErrorSuite(err error) {
	var its suite.Suite
	if err != nil {
		its.FailNow(err.Error())
	}
}
