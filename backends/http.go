package backends

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/moisespsena-go/logging"
)

var log_ = logging.GetOrCreateLogger("github.com/moisespsena-go/logging/backends")

type HttpOptions struct {
	Timeout   int
	Insecure  bool
	HttpGet   bool
	Formatted bool
	Async     bool
}

type HttpBackend struct {
	Client        *http.Client
	URL           url.URL
	HttpGet       bool
	Formatted     bool
	defaultClient bool
	Async         bool
}

func NewHttpBackend(URL url.URL, opt HttpOptions, client *http.Client) (wsb *HttpBackend) {
	var defaultClient bool
	if client == nil {
		dd := *http.DefaultClient
		client = &dd
		defaultClient = true
	}
	if opt.Timeout == 0 {
		opt.Timeout = 2
	}

	client.Timeout = time.Second * time.Duration(opt.Timeout)
	if client.Transport == nil {
		transport := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: client.Timeout,
			}).DialContext,
			TLSHandshakeTimeout: 2 * time.Second,
		}
		if opt.Insecure {
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
		client.Transport = transport
	}

	wsb = &HttpBackend{
		Client:        client,
		URL:           URL,
		HttpGet:       opt.HttpGet,
		Formatted:     opt.Formatted,
		defaultClient: defaultClient,
		Async:         opt.Async,
	}
	return
}

func (this HttpBackend) log(level logging.Level, calldepth int, rec *logging.Record) (err error) {
	var msg []byte
	if this.Formatted {
		msg = []byte(rec.Formatted(calldepth))
	} else if msg, err = json.Marshal(rec.Data()); err != nil {
		return
	}
	var resp *http.Response
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if this.HttpGet {
		var url = this.URL
		url.Query().Set("message", string(msg))
		_, err = this.Client.Get(url.String())
	} else {
		_, err = this.Client.Post(this.URL.String(), "application/json", bytes.NewBuffer(msg))
	}
	return
}

func (this HttpBackend) print(args ...interface{}) (err error) {
	msg := []byte(fmt.Sprint(args...))
	var resp *http.Response
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if this.HttpGet {
		var url = this.URL
		url.Query().Set("string", string(msg))
		_, err = this.Client.Get(url.String())
	} else {
		var url = this.URL
		url.Query().Set("string", "true")
		_, err = this.Client.Post(url.String(), "application/json", bytes.NewBuffer(msg))
	}
	return
}

func (this HttpBackend) Print(args ...interface{}) (err error) {
	if this.Async {
		go func() {
			if err := this.print(args...); err != nil {
				log_.Errorf("http async %q failed: %s", this.URL.String(), err.Error())
			}
		}()
	} else {
		err = this.print(args...)
	}
	return
}

func (this HttpBackend) Log(level logging.Level, calldepth int, rec *logging.Record) (err error) {
	if this.Async {
		go func() {
			r := *rec
			if err := this.log(level, calldepth, &r); err != nil {
				log_.Errorf("http async %q failed: %s", this.URL.String(), err.Error())
			}
		}()
	} else {
		err = this.log(level, calldepth, rec)
	}
	return
}

func (this HttpBackend) Close() error {
	if !this.defaultClient {
		this.Client.CloseIdleConnections()
	}
	return nil
}
