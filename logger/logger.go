package logger

import (
	"github.com/op/go-logging"
	"github.com/sethdmoore/digo/types"
	"os"
	//"strings"
)

func Init(c *types.Config) *logging.Logger {
	var log = logging.MustGetLogger("Digo")

	var stdout_format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} %{shortfunc} > %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)

	backend1Formatter := logging.NewBackendFormatter(backend1, stdout_format)
	//backend1Level := logging.AddModuleLevel(backend1)

	logging.SetBackend(backend1Formatter)
	log.Debug("Logger initialized")

	return log
}
