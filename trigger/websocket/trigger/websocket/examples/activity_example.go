package examples

import (
	trigger "github.com/project-flogo/contrib/trigger/websocket/trigger/websocket"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/microgateway"
	"github.com/project-flogo/contrib/trigger/websocket/activity/wsproxy"
	microapi "github.com/project-flogo/microgateway/api"
)

// Example returns an API example
func Example() (engine.Engine, error) {
	app := api.NewApp()

	gateway := microapi.New("WSProxy")

	serviceWS := gateway.NewService("WSProxy", &wsproxy.Activity{})
	serviceWS.SetDescription("Websocket Activity service")
	serviceWS.AddSetting("uri", "ws://localhost:8080/ws")
	serviceWS.AddSetting("maxconnections", "2")

	step := gateway.NewStep(serviceWS)
	step.AddInput("wsconnection", "=$.payload.wsconnection")

	settings, err := gateway.AddResource(app)
	if err != nil {
		return nil, err
	}

	trg := app.NewTrigger(&trigger.Trigger{}, &trigger.Settings{
		Port: 9096,
		EnabledTLS: false,
		ServerCert: "",
		ServerKey: "",
		ClientAuthEnabled: false,
		TrustStore: "",
	})
	handler, err := trg.NewHandler(&trigger.HandlerSettings{
		Method: "GET",
		Path: "/ws",
		Mode: "2",
	})
	if err != nil {
		return nil, err
	}

	_, err = handler.NewAction(&microgateway.Action{}, settings)
	if err != nil {
		return nil, err
	}

	return api.NewEngine(app)
}

