package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/akaritrading/libs/db"
	"github.com/akaritrading/libs/util"
)

const sessionTokenHeader = "X-Session-Token"

var sessiontoken string
var user db.User

var client = &http.Client{
	Timeout: time.Second * 5,
}

func login() error {

	creds := db.Credential{
		Email:    util.GetEnvVar("USER_EMAIL", "vaxkbihm@sharklasers.com"),
		Password: util.GetEnvVar("USER_PASSWORD", "password"),
	}

	data, _ := json.Marshal(creds)

	newRequest("POST", "/api/login", data)

	res := newRequest("POST", "/api/login", data)
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&user)
	sessiontoken = res.Header.Get(sessionTokenHeader)

	return nil
}

func newRequest(method string, url string, body []byte) *http.Request {

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("X-Session-Token", sessiontoken)

	return req
}
