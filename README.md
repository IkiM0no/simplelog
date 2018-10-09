## `simplelog`

Package simplelog implements a simple logging package for json log events with variadic functional options for messages and event objects.

### Usage  
```
package main

import (
        . "github.com/IkiM0no/simplelog"
)

func main() {
        l, err := NewLogger("foo_app")
        if err != nil {
                panic(err)
        }

        mockEvent := map[string]interface{}{
                "foo":    4,
                "baz":    "some log message",
                "float1": 1.22323,
                "float2": 99.39292,
                "list":   []string{"one", "two", "three"},
                "map":    map[string]int{"one": 1, "two": 2},
        }

        l.Info(Msg("some info message"), Event(mockEvent))
        l.Debug(Msg("and then a debug message"), Event(mockEvent))
        l.Warn(Msg("warning!"), Event(mockEvent))
        l.Error(Msg("error message!"), Event(mockEvent))

        l.Debug(Msg("debug, no event"))
        l.Info(Event(mockEvent))

        var e = Error("there was an error while foo-ing")
        l.Info(MsgF("failed to foo. err: %v", e))

        l.Fatal(Msg("h0z3d!!!"), Event(mockEvent))
}

type Error string

func (e Error) Error() string { return string(e) }
```