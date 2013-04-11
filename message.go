package pusher

/*
 * message handling
 */
type Message struct {
	Name      string   `json:"name"`                //Event name (required)
	Data      string   `json:"data"`                //Event data (required) - limited to 10KB
	Channels  []string `json:"channels,omitempty"`  // Array of one or more channel names - limited to 100 channels
	Channel   string   `json:"channel,omitempty"`   //Channel name if publishing to a single channel (can be used instead of channels)
	Socket_id string   `json:"socket_id,omitempty"` //Excludes the event from being sent to a specific connection (see excluding recipients)
}

func NewMessage(name, data string) *Message {
	return &Message{name, data, []string{}, "", ""}
}
