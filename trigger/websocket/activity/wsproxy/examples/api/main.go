package main

import (
	"os"
	"github.com/project-flogo/core/engine"
	"github.com/AkshayGadikar/contrib/trigger/websocket/activity/wsproxy/examples"
)

func main() {

	e, err := examples.Example(os.Args[1])
	if err != nil {
		panic(err)
	}
	engine.RunEngine(e)
}