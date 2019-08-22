package backends

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/moisespsena-go/logging"
)

var fileMap sync.Map

type FileOptions struct {
	Async    bool
	Truncate bool
	Perm     os.FileMode
}

type WriteCloserBackend struct {
	io.Closer
	logging.Backend
	Name  string
	Async bool
}

func (this WriteCloserBackend) Log(level logging.Level, calldepth int, rec *logging.Record) (err error) {
	if this.Async {
		go func() {
			r := *rec
			if err := this.Backend.Log(level, calldepth, &r); err != nil {
				log_.Errorf("http async %q failed: %s", this.Name, err.Error())
			}
		}()
		return
	}
	return this.Backend.Log(level, calldepth, rec)
}

func (this WriteCloserBackend) Close() error {
	if this.Closer != nil {
		return this.Closer.Close()
	}
	return nil
}

func NewFileBackend(path string, options FileOptions) (b *FilePrintBackend, err error) {
	var f *os.File
	if options.Perm == 0 {
		options.Perm = 0666
	}

	if v, ok := fileMap.Load(path); ok {
		b = v.(*FilePrintBackend)
		return
	}

	if options.Truncate {
		f, err = os.Create(path)
	} else {
		f, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, options.Perm)
	}
	if err != nil {
		return
	}

	b = &FilePrintBackend{
		&WriteCloserBackend{
			Closer:  f,
			Name:    "file:" + path,
			Backend: logging.NewLogBackend(f, "", log.LstdFlags),
			Async:   options.Async,
		}, func(args ...interface{}) (err error) {
			_, err = f.WriteString(fmt.Sprint(args...)+"\n")
			return
		},
	}
	fileMap.Store(path, b)
	return
}

type FilePrintBackend struct {
	*WriteCloserBackend
	PrintFunc func(args ...interface{}) (err error)
}

func (this FilePrintBackend) Print(args ...interface{}) (err error) {
	return this.PrintFunc(args...)
}
