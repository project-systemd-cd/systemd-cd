package logger

import (
	"fmt"
	"reflect"
)

// singleton variable
var logger *LoggerI

func Init(l LoggerI) *LoggerI {
	logger = &l
	return logger
}

func Logger() LoggerI {
	if logger == nil {
		panic("domain/logger not initialized")
	}
	return *logger
}

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type LoggerI interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	SetLevel(level Level) error
}

type Var struct {
	Name  string
	Value any
}

func Var2Text(text string, variables []Var) string {
	if len(variables) != 0 && text != "" {
		text += ":"
	}
	for _, v := range variables {
		var name string
		if v.Name != "" {
			name = v.Name
		} else {
			name = reflect.TypeOf(v.Value).String()
		}
		text += fmt.Sprintf("\n\t%s: %+v", name, v.Value)
	}
	return text
}
