package websocket

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"github.com/project-flogo/core/trigger"
	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/data/metadata"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)


var triggerMd = trigger.NewMetadata(&Settings{}, &Output{}, &HandlerSettings{})

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
	server   *Server
	runner   action.Runner
	wsconn   *websocket.Conn
	settings *Settings
	logger   log.Logger
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{settings: s}, nil
}

//Initialize
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	t.logger = ctx.Logger()
	router := httprouter.New()
	addr := ":" + strconv.Itoa(t.settings.Port)

	if t.settings == nil {
		panic(fmt.Sprintf("No Settings found for trigger"))
	}
	//Check whether TLS (Transport Layer Security) is enabled for the trigger
	enableTLS := false
	serverCert := ""
	serverKey := ""
	if t.settings.EnabledTLS != false {
		enableTLSSetting:= t.settings.EnabledTLS
		if enableTLSSetting {
			//TLS is enabled, get server certificate & key
			enableTLS = true
			if t.settings.ServerCert == "" {
				panic(fmt.Sprintf("No serverCert found for trigger in settings"))
			}
			serverCert = t.settings.ServerCert

			if t.settings.ServerKey == "" {
				panic(fmt.Sprintf("No serverKey found for trigger in settings"))
			}
			serverKey = t.settings.ServerKey
		}
	}
	//Check whether client auth is enabled
	enableClientAuth := false
	trustStore := ""
	if t.settings.ClientAuthEnabled != false {
		enableClientAuthSetting:= t.settings.ClientAuthEnabled
		if enableClientAuthSetting {
			enableClientAuth = true
			if t.settings.TrustStore == "" {
				panic(fmt.Sprintf("Client auth is enabled but client trust store is not provided for trigger in settings"))
			}
			trustStore = t.settings.TrustStore
		}
	}

	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}

		method := s.Method
		path := s.Path
		mode := s.Mode

		router.Handle(method, path, newActionHandler(t, handler,mode))
	}

	t.logger.Debugf("Configured on port %d", t.settings.Port)
	t.server = NewServer(addr, router, enableTLS, serverCert, serverKey, enableClientAuth, trustStore)

	return nil
}

func (t *Trigger) Start() error {
	return t.server.Start()
}


func (t *Trigger) Stop() error {
	t.wsconn.Close()
	return t.server.Stop()
}

func newActionHandler(rt *Trigger, handler trigger.Handler, mode string) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Println("received incomming request")

		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		rt.wsconn = conn
		if err != nil {
			fmt.Println("upgrade error", err)
		} else {
			defer conn.Close()
			//upgraded to websocket connection
			clientAdd := conn.RemoteAddr()
			fmt.Println("Upgraded to websocket protocol")
			fmt.Println("Remote address:", clientAdd)
			if mode == "1" {
				for {
					_, message, err := rt.wsconn.ReadMessage()
					fmt.Println("Message received :", string(message))
					if err != nil {
						fmt.Errorf("error while reading websocket message: %s", err)
						break
					}
					out := &Output{}
					out.Content = message
					_, err = handler.Handle(context.Background(), out)
					if err != nil {
						fmt.Errorf("Run action  failed [%s] ", err)
					}

				}
				rt.logger.Infof("stopped listening to websocket endpoint")
			}
			if mode == "2" {
				out := &Output{}
				out.QueryParams = make(map[string]string)
				out.PathParams = make(map[string]string)
				out.Headers = make(map[string]string)
				out.WSconnection = conn
				_, err := handler.Handle(context.Background(), out)
				if err != nil {
					fmt.Errorf("Run action  failed [%s] ", err)
				}
			}
		}

	}
}