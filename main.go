package main

import (
	"flag"

	"github.com/KonishchevDmitry/rsspipes"
	"github.com/KonishchevDmitry/rsspipes/util"

	"fmt"
	_ "github.com/KonishchevDmitry/my-rsspipes/pipes"
)

var log = util.MustGetLogger("server")

func main() {
	util.MustInitLogging(false, false)

	portArg := flag.Int("port", 8003, "bind port")
	flag.Parse()

	err := rsspipes.Serve(fmt.Sprintf("localhost:%d", *portArg))
	if err != nil {
		log.Fatal(err)
	}
}
