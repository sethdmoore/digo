package logger

/*
	Given the sheer amount of work it is to set up a fucking logger
	I should probably re-use this package elsewhere
*/

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/types"
	"os"
	"strings"
)

func set_loglevel(input string, c *types.Config) logging.Level {
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

func open_logfile(log_location string) (*os.File, bool) {
	var f *os.File
	var err error
	_, err = os.Stat(log_location)
	if err != nil {
		fmt.Printf("WARN: log file could not be read: %s\n", err)
		f, err = os.Create(log_location)
		if err == nil {
			return f, true
		} else {
			fmt.Printf("WARN: could not create log file: %s\n", err)
			return nil, false
		}
	}

	f, err = os.OpenFile(log_location, os.SEEK_END, os.ModeAppend)
	if err == nil {
		return f, true
	} else {
		fmt.Printf("Could not open log file for writing: %s\n", err)
		return nil, false
	}
}

func Init(c *types.Config) *logging.Logger {
	var log = logging.MustGetLogger("Digo")
	var llevel logging.Level
	var use_stdout bool
	var use_logfile bool

	// scoping made this a requirement //
	var logfile_backend *logging.LogBackend
	var logfile_backendFormatter logging.Backend
	var logfile_backend_leveled logging.LeveledBackend

	var stdout_backend *logging.LogBackend
	var stdout_backendFormatter logging.Backend
	var stdout_backend_leveled logging.LeveledBackend
	//                                 //

	// log streams. ["stdout"] || ["file"] || ["stdout", "file"]
	for _, item := range strings.Split(c.LogStreams, ",") {
		switch {
		case item == "stdout":
			use_stdout = true
		case item == "file":
			use_logfile = true
		}
	}

	llevel = set_loglevel(c.LogLevel, c)

	// If you want to completely suppress program output, redirect it to /dev/null
	if !use_stdout && !use_logfile {
		fmt.Println("Explicitly enabling stdout")
		use_stdout = true
		use_logfile = false // probably don't explicitly need this
	}

	if use_logfile {
		logfile, success := open_logfile(c.LogFile)
		if success {
			var logfile_format = logging.MustStringFormatter(
				`%{time:15:04:05} %{shortfunc} > %{level:.4s} %{id:03x} %{message}`,
			)
			logfile_backend = logging.NewLogBackend(logfile, "", 0)
			logfile_backendFormatter = logging.NewBackendFormatter(logfile_backend, logfile_format)
			logfile_backend_leveled = logging.AddModuleLevel(logfile_backendFormatter)
			//logfile_backend_leveled.SetLevel(llevel, "Digo")
			logfile_backend_leveled.SetLevel(logging.DEBUG, "Digo")
			// log stuff
		} else {
			fmt.Println("Disabled logfile")
			use_logfile = false
		}
	}

	if use_stdout {
		var stdout_format = logging.MustStringFormatter(
			`%{color}%{time:15:04:05} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
		)
		stdout_backend = logging.NewLogBackend(os.Stdout, "", 0)
		stdout_backendFormatter = logging.NewBackendFormatter(stdout_backend, stdout_format)
		stdout_backend_leveled = logging.AddModuleLevel(stdout_backendFormatter)
		stdout_backend_leveled.SetLevel(llevel, "Digo")
	}

	switch {
	case use_stdout && use_logfile:
		logging.SetBackend(stdout_backend_leveled, logfile_backend_leveled)
	case use_stdout:
		logging.SetBackend(stdout_backend_leveled)
	case use_logfile:
		logging.SetBackend(logfile_backend_leveled)
	default:
		fmt.Println("Error: Could not enable any output.")
		os.Exit(2)

	}
	//logfile_backend := logging.NewLogBackend(os.Stdout, "", 0)

	//stdout_backendLevel := logging.AddModuleLevel(stdout_backend)

	//stdout_backend_leveled := logging.AddModuleLevel(stdout_backend)

	log.Debug("Logger initialized")

	return log
}
