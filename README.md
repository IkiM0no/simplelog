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
Output:
```
{"time":"2018-10-09T03:46:06.717Z","uuid":"575ef4eb-ca12-4492-903e-6361f9ee702c","host":"<host>","app":"foo_app","level":"INFO","msg":"some info message","event":{"baz":"some log message","float1":1.22323,"float2":99.39292,"foo":4,"list":["one","two","three"],"map":{"one":1,"two":2}}}
{"time":"2018-10-09T03:46:06.717Z","uuid":"b408714a-02a7-4744-84a6-29aba1a26d8a","host":"<host>","app":"foo_app","level":"DEBUG","msg":"and then a debug message","event":{"baz":"some log message","float1":1.22323,"float2":99.39292,"foo":4,"list":["one","two","three"],"map":{"one":1,"two":2}}}
{"time":"2018-10-09T03:46:06.717Z","uuid":"e900871e-2090-4f9a-9dfb-6fc3fbdc631d","host":"<host>","app":"foo_app","level":"WARN","msg":"warning!","event":{"baz":"some log message","float1":1.22323,"float2":99.39292,"foo":4,"list":["one","two","three"],"map":{"one":1,"two":2}}}
{"time":"2018-10-09T03:46:06.717Z","uuid":"5fc92887-6a0e-4880-904b-a4035fd2b84c","host":"<host>","app":"foo_app","level":"ERROR","msg":"error message!
","event":{"baz":"some log message","float1":1.22323,"float2":99.39292,"foo":4,"list":["one","two","three"],"map":{"one":1,"two":2}}}
{"time":"2018-10-09T03:46:06.717Z","uuid":"981b35ec-c137-44d4-8f5d-82e043fee8c5","host":"<host>","app":"foo_app","level":"DEBUG","msg":"debug, no even
t"}{"time":"2018-10-09T03:46:06.717Z","uuid":"2cf7ca79-12f5-4042-88b6-664c63a2ee7e","host":"<host>","app":"foo_app","level":"INFO","event":{"baz":"some l
og message","float1":1.22323,"float2":99.39292,"foo":4,"list":["one","two","three"],"map":{"one":1,"two":2}}}
{"time":"2018-10-09T03:46:06.718Z","uuid":"def7f3af-2b57-4cb0-9163-8d0057ed1608","host":"<host>","app":"foo_app","level":"INFO","msg":"failed to foo.
err: there was an error while foo-ing"}
{"time":"2018-10-09T03:46:06.718Z","uuid":"e358b047-4ae8-4001-88d6-b95ccd14ed06","host":"<host>","app":"foo_app","level":"FATAL","msg":"h0z3d!!!","event":{"baz":"some log message","float1":1.22323,"float2":99.39292,"foo":4,"list":["one","two","three"],"map":{"one":1,"two":2}}}exit status 2
```