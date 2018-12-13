package main

import (
	"github.com/project-flogo/core/engine"
	"github.com/AkshayGadikar/contrib/trigger/websocket/trigger/websocket/examples"
)

func main() {

	e, err := examples.Example("2", "3")
	if err != nil {
		panic(err)
	}
	engine.RunEngine(e)
}