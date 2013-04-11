package pusher

import (
	"testing"
	"time"
)

var auth = &Auth{
	"",
	"278d425bdf160c739803",
	"7ad3773142a6692b25b8",
}

func Test_sign(t *testing.T) {
	var data = []struct {
		request *Request
		auth    *Auth
		want    string
	}{
		{
			&Request{
				"POST",
				"/apps/3/events",
				1353088179,
				"ec365a775a4cd0599faeb73354201b6f",
			},
			auth,
			"da454824c97ba181a32ccc17a72625ba02771f50b50e1e7430e47a1f3f457e6c",
		},
	}

	for k, v := range data {
		sig := v.request.sign(v.auth)
		if v.want != sig {
			t.Errorf("test %d\n\texpected %s, got %s", k, v.want, sig)
		}
	}
}

func Test_NewRequest(t *testing.T) {
	var data = []struct {
		method   string
		path     string
		body     []byte
		body_md5 string
	}{
		{
			"",
			"",
			[]byte("{\"name\":\"foo\",\"channels\":[\"project-3\"],\"data\":\"{\\\"some\\\":\\\"data\\\"}\"}"),
			"ec365a775a4cd0599faeb73354201b6f"},
	}

	for k, v := range data {
		r := NewRequest(v.body, v.method, v.path)

		now := time.Now().Unix()

		if v.body_md5 != r.body_md5 || (r.auth_timestamp-now) > 1 || r.auth_timestamp == 0 {
			t.Errorf("test %d\n\texpected %s, got %s", k, v.body_md5, r.body_md5)
		}
	}
}

func Test_buildUrl(t *testing.T) {
	var data = []struct {
		request        *Request
		auth           *Auth
		auth_timestamp int64
		want           string
	}{
		{
			NewRequest([]byte("{}"), "POST", "/apps/3/events"),
			auth,
			1234,
			"/apps/3/events?auth_key=278d425bdf160c739803&auth_timestamp=1234&auth_version=1.0&body_md5=99914b932bd37a50b983c5e7c90ae93b&auth_signature=dd2e07c60ddf2852c15222c726e6b3411fc1a7789e246e81cda23572db45157d",
		},
	}

	for k, v := range data {
		v.request.auth_timestamp = v.auth_timestamp // set it to a static value

		url := v.request.buildUrl(v.auth)
		if url != v.want {
			t.Errorf("test %d\n\twant %s\n\t got  %s", k, v.want, url)
		}
	}
}
