package main

import (
    "github.com/KonishchevDmitry/rsspipes"
    "github.com/KonishchevDmitry/rsspipes/util"

    _ "github.com/KonishchevDmitry/my-rsspipes/pipes"
)

var log = util.MustGetLogger("server")

func main() {
    util.MustInitLogging(false, false)

    err := rsspipes.Serve(":8003")
    if err != nil {
        log.Fatal(err)
    }
}
