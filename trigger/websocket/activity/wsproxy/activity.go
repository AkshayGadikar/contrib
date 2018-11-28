package wsproxy

import (
	"github.com/gorilla/websocket"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	activity.Register(&Activity{}, New)
}

const (
	defaultMaxConnections = 5
)

// WSProxy is websocket proxy service
type WSProxy struct {
	serviceName    string
	backendURL     string
	maxConnections int
	clientConn     *websocket.Conn
}

var activityMd = activity.ToMetadata(&Settings{})

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	act := &Activity{settings: s}
	return act, nil
}

// Activity is an activity that is used to invoke a Web socket operation
// settings : {wsconnection, url, maxconnections}
type Activity struct {
	settings      *Settings
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Invokes a web socket operation
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	wspService := &WSProxy{
		serviceName:   ctx.Name(),
		clientConn:a.settings.WSconnection,
		backendURL:a.settings.Uri,
	}
	if a.settings.maxConnections == ""{
		wspService.maxConnections = defaultMaxConnections
	}else{
		wspService.maxConnections = a.settings.maxConnections
	}

	// start proxy client as a goroutine
	go startProxyClient(wspService)
	return true, nil
}
