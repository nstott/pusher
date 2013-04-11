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

//publish data with a specic name to a channel
func PublishEvent(name, data, channel string, auth *Auth) ([]byte, error) {
	message := &Message{name, data, nil, channel, ""}
	return sendPost(fmt.Sprintf("/apps/%s/events", auth.app_id), message, auth)
}

//publish data with a specific name to array of channels
func BroadcastEvent(name, data string, channels []string, auth *Auth) ([]byte, error) {
	message := &Message{name, data, channels, "", ""}
	return sendPost(fmt.Sprintf("/apps/%s/events", auth.app_id), message, auth)
}

func sendPost(path string, message *Message, auth *Auth) ([]byte, error) {

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	// generate and sign the request
	r := NewRequest(jsonMessage, "POST", path)
	url := r.buildUrl(auth)

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
