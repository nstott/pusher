package pusher

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var auth2 = &Auth{
	"",
	"41195",
	"1dd19b6bb91e2f5e1a37",
	"0a8c6be51a2b94b04eca",
}

type TestResponse struct {
	Method string
	Path   string
	Body   string
}

func Test_PublishEvent(t *testing.T) {
	dummy := httptest.NewServer(http.HandlerFunc(testHandler))
	defer dummy.Close()
	auth2.endpoint = "http://" + dummy.Listener.Addr().String()

	var data = []struct {
		name    string
		data    string
		channel string
		want    *TestResponse
	}{
		{
			"my_event",
			"Heyo!",
			"test_channel",
			&TestResponse{
				"POST",
				"/apps/41195/events",
				"{\"name\":\"my_event\",\"data\":\"Heyo!\",\"channel\":\"test_channel\"}",
			},
		},
	}

	for k, v := range data {
		resp, err := PublishEvent(v.name, v.data, v.channel, auth2)
		if err != nil {
			t.Error(err)
		}

		testResponse := new(TestResponse)
		err = json.Unmarshal(resp, testResponse)
		if err != nil {
			t.Errorf("Cannot Unmarshal response")
		}

		if !reflect.DeepEqual(testResponse, v.want) {
			t.Errorf("test %d\n\texpected %s, got %s", k, v.want, resp)
		}
	}
}

/*
 * test http handler, build a response including the http method, the path, and the message
 */
func testHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Cannot read body"))
	}

	testResponse := &TestResponse{r.Method, r.RequestURI[:strings.Index(r.RequestURI, "?")], string(body)}

	resp, err := json.Marshal(testResponse)
	if err != nil {
		w.Write([]byte("Cannot encode response data"))
	}

	w.Write([]byte(resp))
}
