package pusher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Auth struct {
	app_id string
	key    string
	secret string
}

func NewAuth(app_id, key, secret string) *Auth {
	return &Auth{app_id, key, secret}
}

type Pusher struct {
	auth *Auth
}

//creates a new pusher struct
func NewPusher(auth *Auth) *Pusher {
	return &Pusher{auth}
}

//publish data with a specic name to a channel
func (p *Pusher) PublishEvent(name, data, channel string) ([]byte, error) {
	message := &Message{name, data, nil, channel, ""}
	return p.sendPost(fmt.Sprintf("/apps/%s/events", p.auth.app_id), message)
}

//publish data with a specific name to array of channels
func (p *Pusher) BroadcastEvent(name, data string, channels []string) ([]byte, error) {
	message := &Message{name, data, channels, "", ""}
	return p.sendPost(fmt.Sprintf("/apps/%s/events", p.auth.app_id), message)
}

func (p *Pusher) sendPost(path string, message *Message) ([]byte, error) {

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	// generate and sign the request
	r := NewRequest(jsonMessage, "POST", path)
	url := r.buildUrl(p.auth)

	resp, err := http.Post(endpoint+url, "application/json", bytes.NewReader(jsonMessage))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // we've opened the pipe, don't forget to close it

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200:
		return body, nil
	case 400:
		return nil, errors.New(string(body))
	case 401:
		return nil, errors.New(string(body))
	case 403:
		return nil, errors.New("Forbidden: Account Disabled or over quota")
	}
	return nil, errors.New("Unknown Error")
}
