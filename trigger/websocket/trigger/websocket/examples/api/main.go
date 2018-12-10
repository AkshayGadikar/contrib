package main

import (
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/contrib/trigger/websocket/trigger/examples"
)

func main() {
	e, err := examples.Example()
	if err != nil {
		panic(err)
	}
	engine.RunEngine(e)
}