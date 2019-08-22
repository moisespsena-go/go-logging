// Copyright 2013, Örjan Persson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logging

import "io"

// defaultBackend is the backend used for all logging calls.
var defaultBackend LeveledBackend

// Backend is the interface which a log backend need to implement to be able to
// be used as a logging backend.
type Backend interface {
	Log(Level, int, *Record) error
}

// BackendCloser is the interface which a closeable log backend need to implement to be able to
// be used as a logging backend.
type BackendCloser interface {
	Backend
	io.Closer
}

type Printer interface {
	Print(args ...interface{}) (err error)
}

type MustPrint func(args ...interface{}) (err error)

func (this MustPrint) Print(args ...interface{}) {
	this(args...)
}

// BackendPrinter is the interface which a log backend with Print method need to implement to be able to
// be used as a logging backend.
type BackendPrinter interface {
	Backend
	Printer
}

// BackendPrintCloser is the interface which a log backend with Print and Closer methods need to implement to be able to
// be used as a logging backend.
type BackendPrintCloser interface {
	BackendPrinter
	io.Closer
}

type backendPrintClose struct {
	BackendPrinter
	io.Closer
}

func NewBackendPrintClose(backend Backend, closer ...io.Closer) BackendCloser {
	var c io.Closer
	for _, c = range closer {
	}
	return &backendClose{Backend: backend, Closer: c}
}

func (this backendPrintClose) Close() error {
	if this.Closer != nil {
		return this.Closer.Close()
	}
	return nil
}

type backendClose struct {
	Backend
	io.Closer
}

func NewBackendClose(backend Backend, closer ...io.Closer) BackendCloser {
	var c io.Closer
	for _, c = range closer {
	}
	return &backendClose{Backend: backend, Closer: c}
}

func (this backendClose) Close() error {
	if this.Closer != nil {
		return this.Closer.Close()
	}
	return nil
}

// SetBackend replaces the backend currently set with the given new logging
// backend.
func SetBackend(backends ...Backend) LeveledBackend {
	var backend Backend
	if len(backends) == 1 {
		backend = backends[0]
	} else {
		backend = MultiLogger(backends...)
	}

	defaultBackend = AddModuleLevel(backend)
	return defaultBackend
}

// SetLevel sets the logging level for the specified module. The module
// corresponds to the string specified in GetOrCreateLogger.
func SetLevel(level Level, module string) {
	defaultBackend.SetLevel(level, module)
}

// GetLevel returns the logging level for the specified module.
func GetLevel(module string) Level {
	return defaultBackend.GetLevel(module)
}

// SetLogLevel sets the logging level for the specified module in Log.
func SetLogLevel(log Logger, level Level, module string) {
	if backend := log.Backend(); backend != nil {
		backend.SetLevel(level, module)
		return
	}
	defaultBackend.SetLevel(level, module)
}

// GetLogLevel returns the logging level for the specified module in Log.
func GetLogLevel(log Logger, module string) Level {
	if backend := log.Backend(); backend != nil {
		return backend.GetLevel(module)
	}
	return defaultBackend.GetLevel(module)
}

func DefaultBackendProxy() LeveledBackend {
	return &LeveledBackendProxy{func() LeveledBackend {
		return defaultBackend
	}}
}

type LeveledBackendProxy struct {
	Get func() LeveledBackend
}

func NewLeveledBackendProxy(get func() LeveledBackend) *LeveledBackendProxy {
	return &LeveledBackendProxy{Get: get}
}

func (this LeveledBackendProxy) Log(level Level, calldepth int, rec *Record) error {
	return this.Get().Log(level, calldepth, rec)
}

func (this LeveledBackendProxy) GetLevel(module string) Level {
	return this.Get().GetLevel(module)
}

func (this LeveledBackendProxy) SetLevel(level Level, module string) {
	this.Get().SetLevel(level, module)
}

func (this LeveledBackendProxy) IsEnabledFor(level Level, module string) bool {
	return this.Get().IsEnabledFor(level, module)
}
