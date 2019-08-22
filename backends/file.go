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
	io.WriteCloser
	logging.Backend
	Name  string
	Async bool
}

func NewWriteCloserBackend(name string, wc io.WriteCloser, async bool) *WriteCloserBackend {
	return &WriteCloserBackend{
		WriteCloser: wc,
		Name:        name,
		Backend:     logging.NewLogBackend(wc, "", log.LstdFlags),
		Async:       async,
	}
}

func (this *WriteCloserBackend) Log(level logging.Level, calldepth int, rec *logging.Record) (err error) {
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

func (this *WriteCloserBackend) Close() error {
	if this.WriteCloser != nil {
		return this.WriteCloser.Close()
	}
	return nil
}

func NewFileBackend(path string, options FileOptions) (b *FileBackend, err error) {
	var f *os.File
	if options.Perm == 0 {
		options.Perm = 0666
	}

	if v, ok := fileMap.Load(path); ok {
		b = v.(*FileBackend)
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

	b = &FileBackend{
		path,
		NewWriteCloserBackend("file:"+path, f, options.Async),
	}
	fileMap.Store(path, b)
	return
}

type FileBackend struct {
	path string
	*WriteCloserBackend
}

func (this *FileBackend) Print(args ...interface{}) (err error) {
	_, err = this.Write([]byte(fmt.Sprint(args...) + "\n"))
	return
}

func (this *FileBackend) Path() string {
	return this.path
}
