package main

import (
    "github.com/KonishchevDmitry/rsspipes"
    "github.com/KonishchevDmitry/rsspipes/util"
)

var log = util.MustGetLogger("server")

func main() {
    util.MustInitLogging(false, true)

    err := rsspipes.Serve(":8003")
    if err != nil {
        log.Fatal(err)
    }
}