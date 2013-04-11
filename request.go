package pusher

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"time"
)

type Request struct {
	method         string
	path           string
	auth_timestamp int64
	body_md5       string
}

var endpoint = "http://api.pusherapp.com"
var auth_version = "1.0"

/*
 * build and sign the request
 */
func NewRequest(body []byte, method, path string) *Request {
	r := &Request{method, path, 0, ""}

	// hash the body
	h := md5.New()
	h.Write(body)
	r.body_md5 = fmt.Sprintf("%x", h.Sum(nil))

	// set the timestamp
	r.auth_timestamp = time.Now().Unix()

	return r
}

func (r *Request) sign(a *Auth) string {
	str := fmt.Sprintf("%s\n%s\nauth_key=%s&auth_timestamp=%d&auth_version=%s&body_md5=%s",
		r.method, r.path, a.key, r.auth_timestamp, auth_version, r.body_md5)

	h := hmac.New(sha256.New, []byte(a.secret))
	h.Write([]byte(str))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (r *Request) buildUrl(a *Auth) string {
	return fmt.Sprintf("%s?auth_key=%s&auth_timestamp=%d&auth_version=%s&body_md5=%s&auth_signature=%s",
		r.path,
		a.key,
		r.auth_timestamp,
		auth_version,
		r.body_md5,
		r.sign(a),
	)
}
