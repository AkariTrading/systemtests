package tests

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akaritrading/libs/db"
)

var ErrorStatusCodeNotOk = errors.New("ErrorStatusCodeNotOk")

func createScript() (*db.Script, error) {

	body := map[string]interface{}{"title": "test"}
	d, _ := json.Marshal(body)

	var script db.Script

	err := NewRequest("POST", PlatformRoute("/api/scripts"), d, &script)

	if err != nil {
		return nil, err
	}

	if script.ID == "" {
		return &script, errors.New("create script - bad ID")
	}

	if script.Title != "test" {
		return &script, errors.New("create script - bad ID")
	}

	return &script, nil
}

func getScripts(userID string) ([]db.Script, error) {

	var scripts []db.Script
	return scripts, NewRequest("GET", PlatformRoute("/api/scripts"), nil, &scripts)
}

func getScript(id string) (*db.Script, error) {
	var script db.Script
	return &script, NewRequest("GET", PlatformRoute(fmt.Sprintf("/api/scripts/%s", id)), nil, &script)
}

func updateScript(id string) (*db.Script, error) {

	var script db.Script

	body := map[string]interface{}{"title": "updated"}
	d, _ := json.Marshal(body)

	err := NewRequest("PUT", PlatformRoute(fmt.Sprintf("/api/scripts/%s", id)), d, &script)
	if err != nil {
		return &script, err
	}

	if script.Title != "updated" {
		return &script, errors.New("title was not updated")
	}

	return &script, nil
}

func deleteScript(id string) error {
	return NewRequest("DELETE", PlatformRoute(fmt.Sprintf("/api/scripts/%s", id)), nil, nil)
}

func (test Test) CRUDScript(t *TestRun) {

	script1, err := createScript()
	if err != nil {
		t.Fail(err)
	}

	script2, err := createScript()
	if err != nil {
		t.Fail(err)
	}

	script, err := getScript(script1.ID)
	if err != nil {
		t.Fail(err)
	}

	if script1.ID != script.ID {
		t.FailStr("create script and get script IDs don't match")
	}

	_, err = updateScript(script1.ID)
	if err != nil {
		t.Fail(err)
	}

	scripts, err := getScripts(script1.ID)
	if err != nil {
		t.Fail(err)
	}

	if len(scripts) < 2 {
		t.Fail(err)
	}

	err = deleteScript(script1.ID)
	if err != nil {
		t.Fail(err)
	}

	err = deleteScript(script2.ID)
	if err != nil {
		t.Fail(err)
	}
}
