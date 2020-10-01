package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/akaritrading/libs/db"
	"github.com/akaritrading/libs/flag"
)

type Test struct{}

type TestRun struct {
	Name   string
	Failed bool
	Logs   []string `json:",omitempty"`
}

func (t *TestRun) FailStr(err string) {
	t.Failed = true
	t.Log(err)
}

func (t *TestRun) Fail(err error) {
	t.Failed = true
	t.Log(err.Error())
}

func (t *TestRun) Log(f string, v ...interface{}) {
	t.Logs = append(t.Logs, fmt.Sprintf(f, v...))
}

var client = &http.Client{
	Timeout: time.Second * 5,
}

const sessionTokenHeader = "X-Session-Token"

var sessiontoken string
var user db.User

func PlatformRoute(r string) string {
	return "http://" + flag.PlatformHost() + r
}

func Login() error {

	creds := db.Credential{
		Email:    flag.GetEnvVar("USER_EMAIL", "jpoqriwy@sharklasers.com"),
		Password: flag.GetEnvVar("USER_PASSWORD", "password"),
	}

	data, _ := json.Marshal(creds)

	req, err := http.NewRequest("POST", PlatformRoute("/auth/login"), bytes.NewReader(data))
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&user)

	sessiontoken = res.Header.Get(sessionTokenHeader)

	return nil
}

func NewRequest(method string, url string, body []byte, ret interface{}) error {

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("X-Session-Token", sessiontoken)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ErrorStatusCodeNotOk
	}

	if ret != nil {
		json.NewDecoder(res.Body).Decode(ret)
	}

	return nil
}
