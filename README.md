## `simplelog`

Package simplelog implements a simple logging package for json log events with variadic functional options for messages and event objects.

### Usage  
```
package main

import (
        . "github.com/IkiM0no/simplelog"
)

func main() {
        log, err := NewLogger("foo_app")
        if err != nil {
                panic(err)
        }

        mockEvent := map[string]interface{}{
                "int":    4,
                "string": "log message",
                "float":  1.22323,
                "list":   []string{"one", "two", "three"},
                "map":    map[string]int{"one": 1, "two": 2},
        }

        // log a simple message
        log.Info(Msg("info message"))

        // log a warning with Event
        log.Warn(Msg("warning!"), Event(mockEvent))

        // log a formatted message with event
        var e = Error("there was an error while foo-ing")
        log.Error(MsgF("failed to foo. err: %v", e), Event(mockEvent))

        // log fatal
        log.Fatal(Msg("h0z3d!!!"), Event(mockEvent))
}

type Error string

func (e Error) Error() string { return string(e) }
```
Output:
```
{"time":"2018-10-09T14:52:00.029Z","uuid":"0728c0bd-c751-4966-85f3-d208d9e1463a","host":"<host>","app":"foo_app","level":"INFO","msg":"info message"}
{"time":"2018-10-09T14:52:00.029Z","uuid":"296e1b80-80a2-47eb-b309-ae5daaa99765","host":"<host>","app":"foo_app","level":"WARN","msg":"warning!","even
t":{"float":1.22323,"int":4,"list":["one","two","three"],"map":{"one":1,"two":2},"string":"log message"}}
{"time":"2018-10-09T14:52:00.030Z","uuid":"fe070e9b-328f-4209-b17c-6b197e5ba013","host":"<host>","app":"foo_app","level":"ERROR","msg":"failed to foo.
 err: there was an error while foo-ing","event":{"float":1.22323,"int":4,"list":["one","two","three"],"map":{"one":1,"two":2},"string":"log message"}}
{"time":"2018-10-09T14:52:00.030Z","uuid":"8e7a3b7a-4092-43d1-bd88-d30160ad8986","host":"<host>","app":"foo_app","level":"FATAL","msg":"h0z3d!!!","eve
nt":{"float":1.22323,"int":4,"list":["one","two","three"],"map":{"one":1,"two":2},"string":"log message"}}exit status 2
```