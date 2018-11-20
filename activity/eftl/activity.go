package eftl

import (
	//"bytes"
	"crypto/tls"
	//"encoding/json"
	//"io"
	"io/ioutil"
	//"net/http"
	//"net/url"
	//"strings"

	"github.com/project-flogo/contrib/activity/eftl/utils"
	"github.com/project-flogo/core/activity"
	//"github.com/project-flogo/core/data/metadata"
)

func init() {
	activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{})

type Activity struct {
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func New(ctx activity.InitContext) (activity.Activity, error) {

	act := &Activity{}
	return act, nil
}

func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	logger := ctx.Logger()
	input := &Input{}
	ctx.GetInputObject(input)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	ca := input.CA
	if ca != "" {
		certificate, err := ioutil.ReadFile(ca)
		if err != nil {
			logger.Errorf("can't open certificate", err)
			return false,err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(certificate)
		tlsConfig = &tls.Config{
			RootCAs: pool,
		}
	}
	id := input.Id
	user := input.User
	password := input.Password
	options := &utils.Options{
		ClientID:  id,
		Username:  user,
		Password:  password,
		TLSConfig: tlsConfig,
	}

	url := input.URL
	errorsChannel := make(chan error, 1)
	connection, err := utils.Connect(url, options, errorsChannel)
	if err != nil {
		logger.Errorf("connection failed: %s", err)
		return false,err
	}
	defer connection.Disconnect()

	content := input.Content
	dest := input.Dest
	if dest != nil {
		err = connection.Publish(utils.Message{
			"_dest":   dest,
			"content": content,
		})
		if err != nil {
			logger.Efforf("failed to publish: %s", err)
			return false, err
		}
	}

	return true, nil
}
