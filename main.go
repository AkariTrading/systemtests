package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akaritrading/libs/db"
	"github.com/akaritrading/libs/util"
)

var DB *db.DB

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		err := login()
		if err != nil {
			log.Fatal(err)
		}

		util.WriteJSON(w, runTests())
	})

	http.ListenAndServe(":8080", nil)
}

func initDB() *db.DB {
	db, err := db.Open(util.PostgresHost(), util.PostgresUser(), util.PostgresDBName(), util.PostgresPassword())
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func runTests() []*TestRun {

	var ret []*TestRun

	ret = append(ret, CRUDScript())
	ret = append(ret, CRUDVersion())
	ret = append(ret, VersionsCountLimited())

	return ret
}

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
