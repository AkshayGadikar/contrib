package main

import (
	"os"
	"github.com/project-flogo/core/engine"
	"github.com/AkshayGadikar/contrib/trigger/websocket/trigger/websocket/examples"
)

func main() {

	e, err := examples.Example(os.Args[1], os.Args[2])
	if err != nil {
		panic(err)
	}
	engine.RunEngine(e)
}