package logger

/*
	Given the sheer amount of work it is to set up a fucking logger
	I should probably re-use this package elsewhere
*/

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/config"
	"github.com/sethdmoore/digo/types"
	"os"
	"strings"
)

func setLoglevel(input string, c *types.Config) logging.Level {
	// set loglevel
	var err error
	var llevel logging.Level
	llevel, err = logging.LogLevel(c.LogLevel)
	if err != nil {
		fmt.Printf("WARN: Incorrect log level: %s\n", c.LogLevel)
		fmt.Println("Valid log levels: debug, info, notice, warning, error, critical")
		fmt.Println("Defaulting to info")
		llevel = logging.INFO
	}
	return llevel
}

func openLogfile(logLocation string) (*os.File, bool) {
	var f *os.File
	var err error
	_, err = os.Stat(logLocation)
	if err != nil {
		fmt.Printf("WARN: log file could not be read: %s\n", err)
		f, err = os.Create(logLocation)
		if err != nil {
			fmt.Printf("WARN: could not create log file: %s\n", err)
			return nil, false
		}

		// blame Gofmt for this extra check. Would be fiiiine to put it in the else case
		if f != nil {
			return f, true
		}
	}

	f, err = os.OpenFile(logLocation, os.SEEK_END, os.ModeAppend)
	if err == nil {
		fmt.Printf("Could not open log file for writing: %s\n", err)
		return nil, false
	}
	return f, true
}

// Init sets up the Logger
func Init() *logging.Logger {
	var log = logging.MustGetLogger("Digo")
	var llevel logging.Level
	var useStdout bool
	var useLogfile bool
	c := config.Get()

	// scoping made this a requirement //
	var logfileBackend *logging.LogBackend
	var logfileBackendFormatter logging.Backend
	var logfileBackendLeveled logging.LeveledBackend

	var stdoutBackend *logging.LogBackend
	var stdoutBackendFormatter logging.Backend
	var stdoutBackenedLeveled logging.LeveledBackend
	//                                 //

	// log streams. ["stdout"] || ["file"] || ["stdout", "file"]
	for _, item := range strings.Split(c.LogStreams, ",") {
		switch {
		case item == "stdout":
			useStdout = true
		case item == "file":
			useLogfile = true
		}
	}

	llevel = setLoglevel(c.LogLevel, c)

	// If you want to completely suppress program output, redirect it to /dev/null
	if !useStdout && !useLogfile {
		fmt.Println("Explicitly enabling stdout")
		useStdout = true
		useLogfile = false // probably don't explicitly need this
	}

	if useLogfile {
		logfile, success := openLogfile(c.LogFile)
		if success {
			var logfileFormat = logging.MustStringFormatter(
				`%{time:15:04:05} %{shortfunc} > %{level:.4s} %{id:03x} %{message}`,
			)
			logfileBackend = logging.NewLogBackend(logfile, "", 0)
			logfileBackendFormatter = logging.NewBackendFormatter(logfileBackend, logfileFormat)
			logfileBackendLeveled = logging.AddModuleLevel(logfileBackendFormatter)
			//logfileBackendLeveled.SetLevel(llevel, "Digo")
			logfileBackendLeveled.SetLevel(logging.DEBUG, "Digo")
			// log stuff
		} else {
			fmt.Println("Disabled logfile")
			useLogfile = false
		}
	}

	if useStdout {
		var stdoutFormat = logging.MustStringFormatter(
			`%{color}%{time:15:04:05} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
		)
		stdoutBackend = logging.NewLogBackend(os.Stdout, "", 0)
		stdoutBackendFormatter = logging.NewBackendFormatter(stdoutBackend, stdoutFormat)
		stdoutBackenedLeveled = logging.AddModuleLevel(stdoutBackendFormatter)
		stdoutBackenedLeveled.SetLevel(llevel, "Digo")
	}

	switch {
	case useStdout && useLogfile:
		logging.SetBackend(stdoutBackenedLeveled, logfileBackendLeveled)
	case useStdout:
		logging.SetBackend(stdoutBackenedLeveled)
	case useLogfile:
		logging.SetBackend(logfileBackendLeveled)
	default:
		fmt.Println("Error: Could not enable any output.")
		os.Exit(2)

	}

	log.Debug("Logger initialized")

	return log
}
