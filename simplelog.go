// Package simplelog implements a simple logging package for json log events
// with variadic functional options for messages and event objects.
package simplelog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"simplelog/utils"
	"time"
)

const ISO_8601 = "2006-01-02T15:04:05.999Z"

type Logger interface {
	Info([]func(*LogEvent))
	Debug([]func(*LogEvent))
	Warn([]func(*LogEvent))
	Error([]func(*LogEvent))
	Fatal([]func(*LogEvent))
}

type LGx struct {
	host string
	app  string
}

type LogEvent struct {
	Time  string                 `json:"time"`
	Uuid  string                 `json:"uuid"`
	Host  string                 `json:"host"`
	App   string                 `json:"app"`
	Level string                 `json:"level"`
	Msg   string                 `json:"msg,omitempty"`
	Event map[string]interface{} `json:"event,omitempty"`
}

// Msg functional option sets LogEvent.Msg to the given string.
func Msg(msg string) func(*LogEvent) {
	return func(e *LogEvent) {
		e.Msg = msg
	}
}

// MsgF functional option sets LogEvent.Msg according to the given
// format and list of arguments.
func MsgF(format string, args ...interface{}) func(*LogEvent) {
	return func(e *LogEvent) {
		e.Msg = fmt.Sprintf(format, args...)
	}
}

// Event functional option sets LogEvent.Event to a map.
func Event(event map[string]interface{}) func(*LogEvent) {
	return func(e *LogEvent) {
		e.Event = event
	}
}

// NewLogger returns a logger with hostname and app intialized.
func NewLogger(app string) (*LGx, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("could not determine hostname %v", err)
	}

	return &LGx{hostname, app}, nil
}

// newEvent returns a LogEvent. Convenience method for log-level functions.
func (l *LGx) newEvent(level string, opts ...func(*LogEvent)) ([]byte, error) {
	u, err := utils.GenerateUUID()
	if err != nil {
		return []byte{}, err
	}
	e := &LogEvent{
		Time:  time.Now().UTC().Format(ISO_8601),
		Uuid:  u,
		Host:  l.host,
		App:   l.app,
		Level: level,
	}

	// apply functional options
	for _, opt := range opts {
		opt(e)
	}

	b, err := json.Marshal(*e)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

const EventErr = "could not generate log event: %v\n"

func (l *LGx) Info(opts ...func(*LogEvent)) {
	b, err := l.newEvent("INFO", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	fmt.Fprintf(os.Stdout, "%+v\n", string(b))
}

func (l *LGx) Debug(opts ...func(*LogEvent)) {
	b, err := l.newEvent("DEBUG", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	fmt.Fprintf(os.Stdout, "%+v\n", string(b))
}

func (l *LGx) Warn(opts ...func(*LogEvent)) {
	b, err := l.newEvent("WARN", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	fmt.Fprintf(os.Stderr, "%+v\n", string(b))
}

func (l *LGx) Error(opts ...func(*LogEvent)) {
	b, err := l.newEvent("ERROR", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	fmt.Fprintf(os.Stderr, "%+v\n", string(b))
}

func (l *LGx) Fatal(opts ...func(*LogEvent)) {
	b, err := l.newEvent("FATAL", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	fmt.Fprintf(os.Stderr, "%+v", string(b))
	os.Exit(2)
}
