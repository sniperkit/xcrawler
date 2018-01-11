package ajii

import (
	"fmt"
	"github.com/op/go-logging"
	"os"
)

type Logger logging.Logger

func GetLogger(name string) *logging.Logger {
	log, _ := logging.GetLogger(fmt.Sprintf("x.%s", name))
	var format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-06 15:04:05.000} [%{module}(%{shortfunc})] ▶ [%{level}] %{id:03x} --%{color:reset} %{message}`,
	)
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(backend)
	logging.SetLevel(logging.DEBUG, "")
	logging.SetFormatter(format)
	return log
}
