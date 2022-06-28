package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// computation
// =========================
// =========================

// Math
// =========================

// max(int, int)
func Max(a, b int) (res int) {
	res = a
	if res < b {
		res = b
	}
	return
}

// relation computation
// =========================

// B in A
func In(dst []any, src any) bool {
	for _, d := range dst {
		if src == d {
			return true
		}
	}
	return false
}

// message
// =========================

// return success msg
func Success(m map[string]any) map[string]any {
	msg := map[string]any{
		"status":  "success",
		"collect": "true",
	}
	for k, v := range m {
		msg[k] = v
	}
	return msg
}

// return error msg
func Error(s string) map[string]any {
	msg := map[string]any{
		"status":  "error",
		"collect": "true",
		"msg":     s,
	}
	return msg
}

// Cookie
// =====================

// Args: (context, key, value, duration_minute)
func SetCookie(c echo.Context, k string, v string, m int) {
	cookie := new(http.Cookie)
	cookie.Name = k
	cookie.Value = v
	cookie.MaxAge = 60 * m
	c.SetCookie(cookie)
}

// Args: (context, key)
func GetCookie(c echo.Context, k string) string {
	v, _ := c.Cookie(k)
	return v.Value
}

// Args: (context)
func RenewCookie(c echo.Context) {

}

// Session
// =====================

// Args: (context, session name, path, maxAge)
func Session(c echo.Context, s string, path string, maxAge int) (*sessions.Session, error) {
	sess, err := session.Get(s, c)
	sess.Options = &sessions.Options{
		Path:     path,
		MaxAge:   maxAge,
		HttpOnly: true,
	}
	return sess, err
}

// Args: (context, session, path, maxAge)
func SetSession(c echo.Context, sess *sessions.Session, kv map[string]any) {
	for k, v := range kv {
		sess.Values[k] = v
	}
	sess.Save(c.Request(), c.Response())
}

// Renew Session
// Args: (context, session name)
func GetSession(c echo.Context, s string) (*sessions.Session, error) {
	sess, err := session.Get(s, c)
	return sess, err
}

// clear session map
// Args: (session)
func ClearSession(sess *sessions.Session) {
	for k := range sess.Values {
		delete(sess.Values, k)
	}
}

// Response
// =====================

func ParseJsonBodyToMap(body io.ReadCloser) map[string]any {
	// parse json body from response to map
	byteRes, _ := ioutil.ReadAll(body)
	fmt.Println(string(byteRes))
	var j any
	json.Unmarshal(byteRes, &j)
	resp_data := j.(map[string]any)
	return resp_data
}

func GenerateMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
