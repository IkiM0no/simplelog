// Package simplelog implements a simple logging package for json log events
// with variadic functional options for messages and event objects.
package simplelog // import github.com/IkiM0no/simplelog

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/IkiM0no/simplelog/flat"
	"github.com/IkiM0no/simplelog/utils"
	"github.com/urfave/negroni"
)

const ISO_8601 = "2006-01-02T15:04:05.999Z"

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	CRITICAL
)

var gDefaultLogLevel = INFO

type Logger interface {
	Print([]func(*LogEvent))

	Tracef(string, []interface{})
	Debugf(string, []interface{})
	Printif(bool, string, []interface{})
	Infof(string, []interface{})
	Warnf(string, []interface{})
	Errorf(string, []interface{})
	Fatalf(string, []interface{})
}

type LGx struct {
	host     string   `json:"host,omitempty"`
	app      string   `json:"app",omitempty`
	mode     string   `json:"mode",omitempty`
	Level    LogLevel `json:"level",omitempty`
	HttpXLog []string
}

type LogEvent struct {
	Event map[string]interface{} `json:"event,omitempty"`
	Time  string                 `json:"time"`
	Uuid  string                 `json:"uuid"`
	Host  string                 `json:"host"`
	App   string                 `json:"app"`
	Level string                 `json:"level"`
	Msg   string                 `json:"msg,omitempty"`
	mode  string                 `json:"mode",omitempty`
}

type Opts struct {
	Event []func(*LogEvent)
}

func lvlFromString(levelStr string) LogLevel {
	levelStr = strings.ToLower(strings.TrimSpace(levelStr))
	switch levelStr {
	case "trace":
		return TRACE
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "critical":
		return CRITICAL
	default:
		return gDefaultLogLevel
	}
}

const kvpTemplate = `"date"="%s" "uuid"="%s" "host"="%s" "app"="%s" "level"="%s" "msg"="%s" %s`

func newEvent(level string, opts ...func(*LogEvent)) ([]byte, error) {
	u, err := utils.GenerateUUID()
	if err != nil {
		return []byte{}, err
	}
	e := &LogEvent{
		Time:  time.Now().UTC().Format(ISO_8601),
		Uuid:  u,
		Level: level,
	}

	for _, opt := range opts {
		opt(e)
	}

	if e.mode == "" {
		e.mode = "kvp"
	}

	if e.mode == "json" {
		b, err := json.Marshal(*e)
		if err != nil {
			return []byte{}, err
		}
		return b, nil
	} else {
		e := *e
		f, err := flat.Flatten(e.Event, "event_")
		if err != nil {
			return []byte{}, err
		}

		s := flat.FlatMap(f)
		event := fmt.Sprintf(kvpTemplate, e.Time, e.Uuid,
			e.Host, e.App, e.Level, e.Msg, s)
		return []byte(event), nil
	}
}

const EventErr = "could not generate log event: %v\n"

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

// Mode functional option sets LogEvent mode
func Mode(mode string) func(*LogEvent) {
	return func(e *LogEvent) {
		e.mode = mode
	}
}

// Host functional option sets LogEvent host
func Host(host string) func(*LogEvent) {
	return func(e *LogEvent) {
		e.Host = host
	}
}

// App functional option sets LogEvent app
func App(app string) func(*LogEvent) {
	return func(e *LogEvent) {
		e.App = app
	}
}

func TraceF(opts ...func(*LogEvent)) {
	b, err := newEvent("TRACE", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	put(string(b), "%s\n", 0)
}

func DebugF(opts ...func(*LogEvent)) {
	b, err := newEvent("DEBUG", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	put(string(b), "%s\n", 1)
}

func InfoF(opts ...func(*LogEvent)) {
	b, err := newEvent("INFO", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	put(string(b), "%s\n", 2)
}

func WarnF(opts ...func(*LogEvent)) {
	b, err := newEvent("WARN", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	put(string(b), "%s\n", 3)
}

func ErrorF(opts ...func(*LogEvent)) {
	b, err := newEvent("ERROR", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	put(string(b), "%s\n", 4)
}

func FatalF(opts ...func(*LogEvent)) {
	b, err := newEvent("CRITICAL", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	put(string(b), "%s\n", 5)
	os.Exit(2)
}

func (l *LGx) put(logEvent, format string, lvl LogLevel) {
	if lvl >= l.Level {
		fmt.Fprintf(os.Stdout, format, logEvent)
	} else {
		fmt.Fprintf(os.Stderr, format, logEvent)
	}
}

func put(logEvent, format string, lvl LogLevel) {
	if lvl >= DEBUG {
		fmt.Fprintf(os.Stdout, format, logEvent)
	} else {
		fmt.Fprintf(os.Stderr, format, logEvent)
	}
}

// SLOG OBJECT FUNCTIONAL METHODS

func (l *LGx) Print(opts ...func(*LogEvent)) {
	b, err := l.newEvent("INFO", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%+v\n", 1)
}

// SLOG FUNCTIONAL METHODS

func (l *LGx) TraceF(opts ...func(*LogEvent)) {
	b, err := l.newEvent("DEBUG", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 0)
}

func (l *LGx) DebugF(opts ...func(*LogEvent)) {
	b, err := l.newEvent("INFO", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 1)
}

func (l *LGx) InfoF(opts ...func(*LogEvent)) {
	b, err := l.newEvent("INFO", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 2)
}

func (l *LGx) InfoiF(print bool, opts ...func(*LogEvent)) {
	if print {
		b, err := l.newEvent("INFO", opts...)
		if err != nil {
			log.Printf(EventErr, err)
			return
		}
		l.put(string(b), "%s\n", 2)
	}
}

func (l *LGx) WarnF(opts ...func(*LogEvent)) {
	b, err := l.newEvent("WARN", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 3)
}

func (l *LGx) ErrorF(opts ...func(*LogEvent)) {
	b, err := l.newEvent("ERROR", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 4)
}

func (l *LGx) FatalF(opts ...func(*LogEvent)) {
	b, err := l.newEvent("CRITICAL", opts...)
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 5)
	os.Exit(2)
}

// SLOG STANDARD METHODS

func (l *LGx) Tracef(format string, v ...interface{}) {
	b, err := l.newEvent("DEBUG", MsgF(format, v...))
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 0)
}

func (l *LGx) Debugf(format string, v ...interface{}) {
	b, err := l.newEvent("INFO", MsgF(format, v...))
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 1)
}

func (l *LGx) Infof(format string, v ...interface{}) {
	b, err := l.newEvent("INFO", MsgF(format, v...))
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 2)
}

func (l *LGx) Printf(format string, v ...interface{}) {
	b, err := l.newEvent("INFO", MsgF(format, v...))
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 2)
}

func (l *LGx) Infoif(print bool, format string, v ...interface{}) {
	if print {
		b, err := l.newEvent("INFO", MsgF(format, v...))
		if err != nil {
			log.Printf(EventErr, err)
			return
		}
		l.put(string(b), "%s\n", 2)
	}
}

func (l *LGx) Warnf(format string, v ...interface{}) {
	b, err := l.newEvent("WARN", MsgF(format, v...))
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 3)
}

func (l *LGx) Errorf(format string, v ...interface{}) {
	b, err := l.newEvent("ERROR", MsgF(format, v...))
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 4)
}

func (l *LGx) Fatalf(format string, v ...interface{}) {
	b, err := l.newEvent("CRITICAL", MsgF(format, v...))
	if err != nil {
		log.Printf(EventErr, err)
		return
	}
	l.put(string(b), "%s\n", 5)
	os.Exit(2)
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

	for _, opt := range opts {
		opt(e)
	}

	if e.mode == "" {
		e.mode = "kvp"
	}

	if l.mode == "json" {
		b, err := json.Marshal(*e)
		if err != nil {
			return []byte{}, err
		}
		return b, nil
	} else {
		e := *e
		f, err := flat.Flatten(e.Event, "event_")
		if err != nil {
			return []byte{}, err
		}

		s := flat.FlatMap(f)
		event := fmt.Sprintf(kvpTemplate, e.Time, e.Uuid,
			e.Host, e.App, e.Level, e.Msg, s)
		return []byte(event), nil
	}
}

// New slog functional options
func WithHost(host string) func(*LGx) {
	return func(l *LGx) {
		l.host = host
	}
}

func WithApp(app string) func(*LGx) {
	return func(l *LGx) {
		l.app = app
	}
}

func WithMode(mode string) func(*LGx) {
	return func(l *LGx) {
		if mode != "json" && mode != "kvp" {
			l.mode = "kvp"
		}
		l.mode = mode
	}
}

func WithLevel(level string) func(*LGx) {
	return func(l *LGx) {
		l.Level = lvlFromString(level)
	}
}

func WithHttpXLog(x []string) func(*LGx) {
	return func(l *LGx) {
		l.HttpXLog = x
	}
}

// New returns a logger with functional options intialized.
func New(opts ...func(*LGx)) *LGx {
	l := &LGx{}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Method ServeHTTP provides portability for passing simplelog to http middleware
func (l *LGx) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now().UTC()
	next(rw, req)
	if !utils.StringInSlice(req.URL.Path, l.HttpXLog) {
		res := rw.(negroni.ResponseWriter)
		logMessage := map[string]interface{}{
			"req_time":      start.Format(ISO_8601),
			"req_status":    strconv.Itoa(res.Status()),
			"req_elapsed":   fmt.Sprintf("%f", time.Since(start).Seconds()*1e3),
			"req_x_fwd_for": req.Header.Get("X-Forwarded-For"),
			"req_method":    req.Method,
			"req_url_path":  req.URL.Path,
		}
		l.Print(Event(logMessage))
	}
}
