package exchange

import (
	"net/url"
	"strings"

	"github.com/apex/log"

	"github.com/mitchellh/mapstructure"

	"github.com/moisespsena-go/logging"
	"github.com/moisespsena-go/logging/backends"
)

type LogLevel struct {
	Level string
}

var levels = map[string]logging.Level{
	"CRITICAL": logging.CRITICAL,
	"C":        logging.CRITICAL,
	"ERROR":    logging.ERROR,
	"E":        logging.ERROR,
	"WARNING":  logging.WARNING,
	"W":        logging.WARNING,
	"NOTICE":   logging.NOTICE,
	"N":        logging.NOTICE,
	"INFO":     logging.INFO,
	"I":        logging.INFO,
	"DEBUG":    logging.DEBUG,
	"D":        logging.DEBUG,
}

func (ll LogLevel) GetLevel(defaul ...logging.Level) logging.Level {
	if l, ok := levels[strings.ToUpper(ll.Level)]; ok {
		return l
	}
	for _, d := range defaul {
		return d
	}
	return logging.DEBUG
}

type ModuleLoggingBackendConfig struct {
	LogLevel `yaml:",inline"`
	Dst      string
	Options  map[string]interface{}
}

type ModuleLoggingConfig struct {
	LogLevel `yaml:",inline"`
	Name     string
	Backends []ModuleLoggingBackendConfig
	Options  map[string]interface{}
}

func (this ModuleLoggingConfig) Backend() (results []logging.BackendCloser) {
	if len(this.Backends) == 0 {
		return
	}

	for i, b := range this.Backends {
		if strings.HasPrefix(b.Dst, "http:") || strings.HasPrefix(b.Dst, "https:") {
			var opts backends.HttpOptions
			opts.Async = true
			err := mapstructure.Decode(b.Options, &opts)
			if err != nil {
				log.Errorf("parse http options for backend #%d `%s` failed: %s", i, b.Dst, err.Error())
				continue
			}
			URL, err := url.Parse(b.Dst)
			if err != nil {
				log.Errorf("parse url for backend #%d `%s` failed: %s", i, b.Dst, err.Error())
				continue
			}
			bce := backends.NewHttpBackend(*URL, opts, nil)
			results = append(results, bce)
		} else if b.Dst == "-" || b.Dst == "_" {
			results = append(results, logging.NewBackendClose(logging.DefaultBackendProxy()))
		} else {
			var opts backends.FileOptions
			opts.Async = true
			err := mapstructure.Decode(b.Options, &opts)
			if err != nil {
				log.Errorf("parse http options for backend #%d `%s` failed: %s", i, b.Dst, err.Error())
				continue
			}
			bce, err := backends.NewFileBackend(b.Dst, opts)
			if err != nil {
				log.Errorf("create file backend for backend #%d `%s` failed: %s", i, b.Dst, err.Error())
				continue
			}
			results = append(results, bce)
		}
	}
	return
}

func (this ModuleLoggingConfig) BackendPrinter() (results []logging.BackendPrintCloser) {
	if len(this.Backends) == 0 {
		return
	}

	for i, b := range this.Backends {
		if strings.HasPrefix(b.Dst, "http:") || strings.HasPrefix(b.Dst, "https:") {
			var opts backends.HttpOptions
			opts.Async = true
			err := mapstructure.Decode(b.Options, &opts)
			if err != nil {
				log.Errorf("parse http options for backend #%d `%s` failed: %s", i, b.Dst, err.Error())
				continue
			}
			URL, err := url.Parse(b.Dst)
			if err != nil {
				log.Errorf("parse url for backend #%d `%s` failed: %s", i, b.Dst, err.Error())
				continue
			}
			bce := backends.NewHttpBackend(*URL, opts, nil)
			results = append(results, bce)
		} else {
			var opts backends.FileOptions
			opts.Async = true
			err := mapstructure.Decode(b.Options, &opts)
			if err != nil {
				log.Errorf("parse http options for backend #%d `%s` failed: %s", i, b.Dst, err.Error())
				continue
			}
			bce, err := backends.NewFileBackend(b.Dst, opts)
			if err != nil {
				log.Errorf("create file backend for backend #%d `%s` failed: %s", i, b.Dst, err.Error())
				continue
			}
			results = append(results, bce)
		}
	}
	return
}

type LoggingConfig struct {
	LogLevel `yaml:",inline"`
	Modules  []ModuleLoggingConfig
}
