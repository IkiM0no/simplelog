## `simplelog`

Package simplelog implements a simple logging package for log events with variadic functional options for messages and event objects.

Key-Value Pair and JSON modes are supported.

### Usage  
```
package main

import (
        slog "learning/slog_dev/simplelog"
)

func main() {

        // Log via the package directly.
        event := map[string]interface{}{
                "foo": 8,
                "bar": "baz",
        }
        slog.InfoF(slog.MsgF("msg: %s", "message!"), slog.Event(event))
        slog.WarnF(slog.MsgF("msg: %s", "message!"), slog.Event(event), slog.Host("my-host"))
        slog.ErrorF(slog.MsgF("msg: %s", "message!"), slog.Event(event), slog.Mode("json"))

        // Or, create a logger with functional options initializing common values.
        sl := slog.New(
                slog.WithHost("my-host"),
                slog.WithMode("json"),
                slog.WithApp("aPPz"),
        )

        sl.Errorf("err: %s", "an error")
        sl.TraceF(slog.Event(event), slog.MsgF("msg: %d", 2))
}

```
Output:
```
"date"="2018-10-31T22:26:18.528Z" "uuid"="15096ea7-d9f8-43d7-847a-089c4704c3cf" "level"="INFO" "msg"="msg: message!" "event_foo"="8" "event_bar"="baz"
"date"="2018-10-31T22:26:18.528Z" "uuid"="61733186-d8a3-4169-ae3d-8a73f4e24673" "host"="my-host" "level"="WARN" "msg"="msg: message!" "event_foo"="8" "event_bar"="baz"
{"event":{"bar":"baz","foo":8},"date":"2018-10-31T22:26:18.528Z","uuid":"997e60cd-e218-4d3a-9fc1-d641ee178be4","level":"ERROR","msg":"msg: message!"}
{"date":"2018-10-31T22:26:18.528Z","uuid":"4c2d2f86-d773-4819-a5a8-5412c9a3cb6a","host":"my-host","app":"aPPz","level":"ERROR","msg":"err: an error"}
{"event":{"bar":"baz","foo":8},"date":"2018-10-31T22:26:18.528Z","uuid":"a3bff9ca-b341-4da8-a4a8-ad2eab10db56","host":"my-host","app":"aPPz","level":"DEBUG","msg":"msg: 2"}
```