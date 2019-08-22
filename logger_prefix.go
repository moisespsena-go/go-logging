package logging

import "strings"

type LogPrefix struct {
	Logger
	prefix string
}

func (this LogPrefix) Parent() Logger {
	return this.Logger
}

func (this LogPrefix) Prefix() string {
	return this.prefix
}

func (this LogPrefix) SetPrefix(v string) {
	this.prefix = v
}

func (this LogPrefix) Fatal(args ...interface{}) {
	this.Logger.Fatal(append([]interface{}{this.prefix}, args...)...)
}

func (this LogPrefix) Fatalf(format string, args ...interface{}) {
	this.Logger.Fatalf(this.prefix+" "+format, args...)
}

func (this LogPrefix) Panic(args ...interface{}) {
	this.Logger.Panic(append([]interface{}{this.prefix}, args...)...)
}

func (this LogPrefix) Panicf(format string, args ...interface{}) {
	this.Logger.Panicf(this.prefix+" "+format, args...)
}

func (this LogPrefix) Critical(args ...interface{}) {
	this.Logger.Critical(append([]interface{}{this.prefix}, args...)...)
}

func (this LogPrefix) Criticalf(format string, args ...interface{}) {
	this.Logger.Criticalf(this.prefix+" "+format, args...)
}

func (this LogPrefix) Error(args ...interface{}) {
	this.Logger.Error(append([]interface{}{this.prefix}, args...)...)
}

func (this LogPrefix) Errorf(format string, args ...interface{}) {
	this.Logger.Errorf(this.prefix+" "+format, args...)
}

func (this LogPrefix) Warning(args ...interface{}) {
	this.Logger.Warning(append([]interface{}{this.prefix}, args...)...)
}

func (this LogPrefix) Warningf(format string, args ...interface{}) {
	this.Logger.Warningf(this.prefix+" "+format, args...)
}

func (this LogPrefix) Notice(args ...interface{}) {
	this.Logger.Notice(append([]interface{}{this.prefix}, args...)...)
}

func (this LogPrefix) Noticef(format string, args ...interface{}) {
	this.Logger.Noticef(this.prefix+" "+format, args...)
}

func (this LogPrefix) Info(args ...interface{}) {
	this.Logger.Info(append([]interface{}{this.prefix}, args...)...)
}

func (this LogPrefix) Infof(format string, args ...interface{}) {
	this.Logger.Infof(this.prefix+" "+format, args...)
}

func (this LogPrefix) Debug(args ...interface{}) {
	this.Logger.Debug(append([]interface{}{this.prefix}, args...)...)
}

func (this LogPrefix) Debugf(format string, args ...interface{}) {
	this.Logger.Debugf(this.prefix+" "+format, args...)
}

func WithPrefix(parent Logger, prefix string, sep ...string) LogPrefixer {
	s := " ->"
	if len(sep) > 0 {
		s = sep[0]
	}
	return &LogPrefix{parent, strings.TrimSpace(prefix) + s}
}
