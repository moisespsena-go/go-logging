package defaultlog

import (
	"os"

	"github.com/moisespsena-go/logging"
)

var (
	Format = logging.MustStringFormatter(
		`%{time:2006-01-02 15:04:05.999 -07:00}%{color} %{pid} %{level:.4s} [%{module}]: %{message}%{color:reset}`,
	)

	GetOrCreateLogger = logging.GetOrCreateLogger
)

func init() {
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), Format))
}
