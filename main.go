package main

import (
	"log"
	"net/http"
	"reflect"

	"github.com/akaritrading/libs/db"
	"github.com/akaritrading/libs/util"
	"github.com/akaritrading/systemtests/tests"
)

var DB *db.DB

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		err := tests.Login()
		if err != nil {
			util.ErrorJSON(w, err)
			return
		}

		util.WriteJSON(w, runTests())
	})

	http.ListenAndServe(":8080", nil)
}

func initDB() *db.DB {
	db, err := db.DefaultOpen()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func runTests() []*tests.TestRun {

	var ret []*tests.TestRun

	testsValue := reflect.ValueOf(tests.Test{})
	testsType := reflect.TypeOf(tests.Test{})

	for i := 0; i < testsValue.NumMethod(); i++ {
		method := testsValue.Method(i)
		run := &tests.TestRun{Name: testsType.Method(i).Name}
		method.Call([]reflect.Value{reflect.ValueOf(run)})
		ret = append(ret, run)
	}

	return ret
}
