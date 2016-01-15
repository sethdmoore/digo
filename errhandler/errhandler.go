package errhandler

import (
	"github.com/op/go-logging"
	"os"
)

/*
	Deprecated. Useless package
*/

var log *logging.Logger

// Handle an error. Really doesn't do anything
func Handle(err error) {
	if err != nil {
		log.Errorf("%v\n", err)
		os.Exit(2)
	}
}

// Init is probably not necessary, set to Get later
func Init(logger *logging.Logger) {
	// set the package var to the ref
	log = logger
}
