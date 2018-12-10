package wsproxy

import (
	"testing"
	//"time"
	"net/http"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	logger "github.com/project-flogo/core/support/log"
	"github.com/stretchr/testify/assert"
	"github.com/gorilla/websocket"
	"github.com/project-flogo/core/data"
)

type initContext struct {
	settings map[string]interface{}
}

func newInitContext(values map[string]interface{}) *initContext {
	if values == nil {
		values = make(map[string]interface{})
	}
	return &initContext{
		settings: values,
	}
}

func (i *initContext) Settings() map[string]interface{} {
	return i.settings
}

func (i *initContext) MapperFactory() mapper.Factory {
	return nil
}

func (i *initContext) Logger() logger.Logger {
	return logger.RootLogger()
}

type activityContext struct {
	input  map[string]interface{}
	output map[string]interface{}
}

func newActivityContext(values map[string]interface{}) *activityContext {
	if values == nil {
		values = make(map[string]interface{})
	}
	return &activityContext{
		input:  values,
	}
}

func (a *activityContext) ActivityHost() activity.Host {
	return a
}

func (a *activityContext) Name() string {
	return "test"
}

func (a *activityContext) GetInput(name string) interface{} {
	return a.input[name]
}

func (a *activityContext) SetOutput(name string, value interface{}) error {
	return nil
}

func (a *activityContext) GetInputObject(input data.StructValue) error {
	return input.FromMap(a.input)
}

func (a *activityContext) SetOutputObject(output data.StructValue) error {
	return nil
}

func (a *activityContext) GetSharedTempData() map[string]interface{} {
	return nil
}

func (a *activityContext) ID() string {
	return "test"
}

func (a *activityContext) IOMetadata() *metadata.IOMetadata {
	return nil
}

func (a *activityContext) Reply(replyData map[string]interface{}, err error) {
}

func (a *activityContext) Return(returnData map[string]interface{}, err error) {
}

func (a *activityContext) Scope() data.Scope {
	return nil
}

func (a *activityContext) Logger() logger.Logger {
	return logger.RootLogger()
}

func TestWSproxy(t *testing.T) {
	fini := make(chan bool, 1)
	wsHandler := func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		for {
			mt, message, err := conn.ReadMessage()
			next := string(message)
			if len(next) > 5 {
				fini <- true
				break
			}
			if err != nil {
				break
			}
			conn.WriteMessage(mt, []byte(next+"."))
			if err != nil {
				break
			}
		}
	}
	middleware := http.NewServeMux()
	middleware.HandleFunc("/ws", wsHandler)
	server := http.Server{
		Addr:    "localhost:8282",
		Handler: middleware,
	}
	done := make(chan bool, 1)
	go func() {
		server.ListenAndServe()
		done <- true
	}()
	defer func() {
		err := server.Shutdown(nil)
		if err != nil {
			t.Fatal(err)
		}
		<-done
	}()

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8282/ws", nil)
	assert.Nil(t, err)
	defer conn.Close()

	activity, err := New(newInitContext(map[string]interface{}{
		"uri": "ws://localhost:8080/ws",
		"maxconnections": "2",
	}))
	assert.Nil(t, err)

	ctx := newActivityContext(map[string]interface{}{
		"wsconnection": conn,
	})
	_, err = activity.Eval(ctx)
	assert.Nil(t, err)

	err = conn.WriteMessage(websocket.TextMessage, []byte("test"))
	assert.Nil(t, err)

	/*select {
	case <-fini:
	case <-time.After(30 * time.Second):
		t.Fatal("test failed: timed out")
	}*/
}

