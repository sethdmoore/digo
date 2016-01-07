package errhandler

import (
	"github.com/op/go-logging"
	"os"
)

/*
	Deprecated. Useless package
*/

var log *logging.Logger

func Handle(err error) {
	if err != nil {
		log.Errorf("%v\n", err)
		os.Exit(2)
	}
}

func Init(logger *logging.Logger) {
	// set the package var to the ref
	log = logger
}
