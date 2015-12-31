package errhandler

import (
	"fmt"
	"os"
)

/*
	Stupid simple error handling
	Probably could and should be more robust
*/

func Handle(err error) {
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(2)
	}
}
