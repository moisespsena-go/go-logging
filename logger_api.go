package logging

// Logger is an interface for types that creates log records based on the functions
// called and passes them to the underlying logging backend.
type Logger interface {
	IsEnabledFor(level Level) bool

	// SetBackend overrides any previously defined backend for this logger.
	SetBackend(backend LeveledBackend)
	// Backend return current backend if has be defined
	Backend() LeveledBackend

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	// Warning logs a message using WARNING as log level.
	Warning(args ...interface{})
	// Warningf logs a message using WARNING as log level.
	Warningf(format string, args ...interface{})
	// Notice logs a message using NOTICE as log level.
	Notice(args ...interface{})
	// Noticef logs a message using NOTICE as log level.
	Noticef(format string, args ...interface{})
	// Info logs a message using INFO as log level.
	Info(args ...interface{})
	// Infof logs a message using INFO as log level.
	Infof(format string, args ...interface{})
	// Debug logs a message using DEBUG as log level.
	Debug(args ...interface{})
	// Debugf logs a message using DEBUG as log level.
	Debugf(format string, args ...interface{})
	// Writer returns the log writer.
	Writer() LogWriter
}

// LogPrefixer is an interface for types that creates log records with prefix.
type LogPrefixer interface {
	Logger
	Prefix() string
	Parent() Logger
}
