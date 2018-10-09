## `simplelog`

Package simplelog implements a simple logging package for json log events with variadic functional options for messages and event objects.

Key-Value Pair and JSON modes are supported.

### Usage  
```
package main

import (
        . "github.com/IkiM0no/simplelog"
)

func main() {
        log, err := NewLogger("foo_app", "json")
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

        kvLog, err := NewLogger("foo_app", "kvp")
        if err != nil {
                panic(err)
        }

        // log a warning with Event in key-value pair format
        kvLog.Warn(Msg("warning!"), Event(mockEvent))
}

type Error string

func (e Error) Error() string { return string(e) }
```
Output:
```
{"time":"2018-10-09T20:59:16.447Z","uuid":"b44730e4-06f2-4ec8-b210-e5bc33ac7dd0","host":"<host>","app":"foo_app","level":"INFO","msg":"info message"}
{"time":"2018-10-09T20:59:16.447Z","uuid":"ad2b2ecc-4c08-410a-9531-75e5156ad30d","host":"<host>","app":"foo_app","level":"WARN","msg":"warning!","even
t":{"float":1.22323,"int":4,"map":{"one":1,"two":2},"string":"log message"}}
{"time":"2018-10-09T20:59:16.447Z","uuid":"7a61e8e4-9b75-4b44-812d-24a166c98824","host":"<host>","app":"foo_app","level":"ERROR","msg":"failed to foo.
 err: there was an error while foo-ing","event":{"float":1.22323,"int":4,"map":{"one":1,"two":2},"string":"log message"}}
 
"time"="2018-10-09T20:59:16.447Z" "uuid"="5449b793-a1ee-41f9-bc6c-7b7e828e9011" "host"="<host>" "app"="foo_app" "level"="WARN" "event_map_one"="1" "ev
ent_map_two"="2" "event_int"="4" "event_string"="log message" "event_float"="1.22323"
```