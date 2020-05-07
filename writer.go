package logging

import "sync/atomic"

type LogWriter interface {
	Write(lvl Level, extraCalldepth int, format *string, args ...interface{})
}

type writerFunc func(lvl Level, extraCalldepth int, format *string, args ...interface{})

func (w writerFunc) Write(lvl Level, extraCalldepth int, format *string, args ...interface{}) {
	w(lvl, extraCalldepth, format, args...)
}

func NewWriter(f func(lvl Level, extraCalldepth int, format *string, args ...interface{})) LogWriter {
	return writerFunc(f)
}

func DefaultWriter(l Logger, module string) LogWriter {
	return NewWriter(func(lvl Level, extraCalldepth int, format *string, args ...interface{}) {
		if !l.IsEnabledFor(lvl) {
			return
		}

		// Create the logging record and pass it in to the backend
		record := &Record{
			ID:     atomic.AddUint64(&sequenceNo, 1),
			Time:   timeNow(),
			Module: module,
			Level:  lvl,
			fmt:    format,
			Args:   args,
		}

		// TODO use channels to fan out the records to all backends?
		// TODO in case of errors, do something (tricky)

		// calldepth=2 brings the stack up to the caller of the level
		// methods, Info(), Fatal(), etc.
		// ExtraCallDepth allows this to be extended further up the stack in case we
		// are wrapping these methods, eg. to expose them package level

		if backend := l.Backend(); backend != nil {
			backend.Log(lvl, 2+extraCalldepth, record)
			return
		}

		defaultBackend.Log(lvl, 2+extraCalldepth, record)
	})
}
