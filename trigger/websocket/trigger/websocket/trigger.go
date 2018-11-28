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
		if err != nil {
			fmt.Println("upgrade error", err)
		} else {
			defer conn.Close()
			//upgraded to websocket connection
			clientAdd := conn.RemoteAddr()
			fmt.Println("Upgraded to websocket protocol")
			fmt.Println("Remote address:", clientAdd)

			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("read error", err)
					break
				}
					if mode == "1" {
						messageToLog := fmt.Sprintf("Received message(%s) from the client", message)
						fmt.Println(messageToLog)
					} else {
						if mode == "2" {
							out := &Output{}
							out.Content = message
							_, err := handler.Handle(context.Background(), out)
							if err != nil {
								fmt.Errorf("Run action  failed [%s] ", err)
							}
						}
					}

			}
			rt.logger.Infof("stopped listening to websocket endpoint")
		}

	}
}


		/*out := &Output{}

		out.PathParams = make(map[string]string)
		for _, param := range ps {
			out.PathParams[param.Key] = param.Value
		}

		queryValues := r.URL.Query()
		out.QueryParams = make(map[string]string, len(queryValues))
		out.Headers = make(map[string]string, len(r.Header))

		for key, value := range r.Header {
			out.Headers[key] = strings.Join(value, ",")
		}

		for key, value := range queryValues {
			out.QueryParams[key] = strings.Join(value, ",")
		}

		// Check the HTTP Header Content-Type
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/x-www-form-urlencoded":
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			s := buf.String()
			m, err := url.ParseQuery(s)
			content := make(map[string]interface{}, 0)
			if err != nil {
				rt.logger.Errorf("Error while parsing query string: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			for key, val := range m {
				if len(val) == 1 {
					content[key] = val[0]
				} else {
					content[key] = val[0]
				}
			}

			out.Content = content
		default:
			var content interface{}
			err := json.NewDecoder(r.Body).Decode(&content)
			if err != nil {
				switch {
				case err == io.EOF:
				// empty body
				case err != nil:
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			out.Content = content
		}

		if (mode == "1"){
			//acts as server....display the details received from client
		}
		results, err := handler.Handle(context.Background(), out)

		reply := &Reply{}
		reply.FromMap(results)

		if err != nil {
			rt.logger.Debugf("Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if reply.Data != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if reply.Code == 0 {
				reply.Code = 200
			}
			w.WriteHeader(reply.Code)
			if err := json.NewEncoder(w).Encode(reply.Data); err != nil {
				log.Error(err)
			}
			return
		}

		if reply.Code > 0 {
			w.WriteHeader(reply.Code)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}*/