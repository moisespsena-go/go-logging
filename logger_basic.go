package logging

import (
	"fmt"
	"os"
)

type Basic struct {
	writer LogWriter

	// ExtraCallDepth can be used to add additional call depth when getting the
	// calling function. This is normally used when wrapping a logger.
	ExtraCalldepth int
}

// NewBasic creates Basic with writer
func NewBasic(writer LogWriter) Basic {
	return Basic{writer: writer}
}

func (l Basic) write(lvl Level, format *string, args ...interface{}) {
	l.writer.Write(lvl, 2+l.ExtraCalldepth, format, args...)
}

// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func (l Basic) Fatal(args ...interface{}) {
	l.write(CRITICAL, nil, args...)
	os.Exit(1)
}

// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
func (l Basic) Fatalf(format string, args ...interface{}) {
	l.write(CRITICAL, &format, args...)
	os.Exit(1)
}

// Panic is equivalent to l.Critical(fmt.Sprint()) followed by a call to panic().
func (l Basic) Panic(args ...interface{}) {
	l.write(CRITICAL, nil, args...)
	panic(fmt.Sprint(args...))
}

// Panicf is equivalent to l.Critical followed by a call to panic().
func (l Basic) Panicf(format string, args ...interface{}) {
	l.write(CRITICAL, &format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Critical logs a message using CRITICAL as log level.
func (l Basic) Critical(args ...interface{}) {
	l.write(CRITICAL, nil, args...)
}

// Criticalf logs a message using CRITICAL as log level.
func (l Basic) Criticalf(format string, args ...interface{}) {
	l.write(CRITICAL, &format, args...)
}

// Error logs a message using ERROR as log level.
func (l Basic) Error(args ...interface{}) {
	l.write(ERROR, nil, args...)
}

// Errorf logs a message using ERROR as log level.
func (l Basic) Errorf(format string, args ...interface{}) {
	l.write(ERROR, &format, args...)
}

// Warning logs a message using WARNING as log level.
func (l Basic) Warning(args ...interface{}) {
	l.write(WARNING, nil, args...)
}

// Warningf logs a message using WARNING as log level.
func (l Basic) Warningf(format string, args ...interface{}) {
	l.write(WARNING, &format, args...)
}

// Notice logs a message using NOTICE as log level.
func (l Basic) Notice(args ...interface{}) {
	l.write(NOTICE, nil, args...)
}

// Noticef logs a message using NOTICE as log level.
func (l Basic) Noticef(format string, args ...interface{}) {
	l.write(NOTICE, &format, args...)
}

// Info logs a message using INFO as log level.
func (l Basic) Info(args ...interface{}) {
	l.write(INFO, nil, args...)
}

// Infof logs a message using INFO as log level.
func (l Basic) Infof(format string, args ...interface{}) {
	l.write(INFO, &format, args...)
}

// Debug logs a message using DEBUG as log level.
func (l Basic) Debug(args ...interface{}) {
	l.write(DEBUG, nil, args...)
}

// Debugf logs a message using DEBUG as log level.
func (l Basic) Debugf(format string, args ...interface{}) {
	l.write(DEBUG, &format, args...)
}

func (l Basic) Writer() LogWriter {
	return l.writer
}
