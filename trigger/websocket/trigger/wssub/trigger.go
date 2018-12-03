package wssub

import (
	"context"
	"fmt"

	"github.com/project-flogo/core/trigger"
	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/data/metadata"
	"github.com/gorilla/websocket"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &Output{})

func init() {
	trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Trigger REST trigger struct
type Trigger struct {
	runner   action.Runner
	wsconn   *websocket.Conn
	settings *Settings
	logger   log.Logger
	config   *trigger.Config
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{settings: s, config:config}, nil
}

//Initialize
func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	t.logger = ctx.Logger()
	urlSetting := t.config.Settings["url"]
	if urlSetting == nil || urlSetting.(string) == "" {
		return fmt.Errorf("server url not provided")
	}

	url := urlSetting.(string)
	t.logger.Infof("dialing websocket endpoint[%s]...", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("error while connecting to websocket endpoint[%s] - %s", url, err)
	}
	t.wsconn = conn
	go func() {
		for {
			_, message, err := t.wsconn.ReadMessage()
			fmt.Println("Message received :", string(message))
			if err != nil {
				fmt.Errorf("error while reading websocket message: %s", err)
				break
			}

			for _, handler := range ctx.GetHandlers() {
				out := &Output{}
				out.Content = message
				_, err := handler.Handle(context.Background(), out)
				if err != nil {
					fmt.Errorf("Run action  failed [%s] ", err)
				}
			}
		}
		t.logger.Infof("stopped listening to websocket endpoint")
	}()
	return nil
}


func (t *Trigger) Start() error {
	return nil
}


func (t *Trigger) Stop() error {
	t.wsconn.Close()
	return nil
}