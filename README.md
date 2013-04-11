a simple library to broadcast messages on pusherapp.com

installation
============
  go get https://github.com/nstott/pusher


usage
=====

once you set up auth,
```
auth := pusher.NewAuth(app_id, key, secret)
```
there are only two interesting end points,
```
func PublishEvent(name, data, channel string, auth *Auth) ([]byte, error) {
```
and
```
func BroadcastEvent(name, data string, channels []string, auth *Auth) ([]byte, error) {
```

example
=======
```
  import (
    "github.com/nstott/pusher"
  )

  func main () {
    auth := pusher.NewAuth(app_id, key, secret)

    _, err := pusher.PublishEvent("my_event", "Heyo!", "test_channel", auth)
    if err != nil {
      // something bad happened
    }
  }
```