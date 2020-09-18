package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/akaritrading/libs/db"
	"github.com/akaritrading/libs/util"
)

var ErrorStatusCodeNotOk = errors.New("ErrorStatusCodeNotOk")

func platformRoute(r string) string {
	return "http://" + util.PlatformHost() + r
}

func createScript(t *TestRun) (*db.Script, error) {

	var script db.Script
	d, _ := json.Marshal(db.Script{Title: "test"})

	req := newRequest("POST", platformRoute("/api/scripts"), d)

	res, err := client.Do(req)
	if err != nil {
		return &script, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&script)

	if script.ID == "" {
		return &script, errors.New("create script - bad ID")
	}

	if res.StatusCode != http.StatusOK {
		return &script, ErrorStatusCodeNotOk
	}

	if script.Title != "test" {
		return &script, errors.New("create script - bad ID")
	}

	return &script, nil
}

func getScripts(t *TestRun, userID string) ([]db.Script, error) {

	route := "/api/scripts/"
	req := newRequest("GET", platformRoute(route), nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var scripts []db.Script
	json.NewDecoder(res.Body).Decode(&scripts)

	if res.StatusCode != http.StatusOK {
		return scripts, ErrorStatusCodeNotOk
	}

	return scripts, nil
}

func getScript(t *TestRun, id string) (*db.Script, error) {
	var script db.Script

	route := fmt.Sprintf("/api/scripts/%s", id)
	req := newRequest("GET", platformRoute(route), nil)
	res, err := client.Do(req)
	if err != nil {
		return &script, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&script)

	if res.StatusCode != http.StatusOK {
		return &script, ErrorStatusCodeNotOk
	}

	return &script, nil
}

func updateScript(t *TestRun, id string) (*db.Script, error) {

	var script db.Script

	d, _ := json.Marshal(db.Script{Title: "updated"})

	route := fmt.Sprintf("/api/scripts/%s", id)
	req := newRequest("PUT", platformRoute(route), d)
	res, err := client.Do(req)
	if err != nil {
		return &script, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&script)

	if res.StatusCode != http.StatusOK {
		return &script, ErrorStatusCodeNotOk
	}

	if script.Title != "updated" {
		return &script, errors.New("title was not updated")
	}

	return &script, nil
}

func deleteScript(t *TestRun, id string) error {

	route := fmt.Sprintf("/api/scripts/%s", id)
	req := newRequest("DELETE", platformRoute(route), nil)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ErrorStatusCodeNotOk
	}

	return nil
}

func CRUDScript() *TestRun {

	t := &TestRun{Name: "CRUDScript"}

	script1, err := createScript(t)
	if err != nil {
		t.Fail(err)
	}

	script2, err := createScript(t)
	if err != nil {
		t.Fail(err)
	}

	script, err := getScript(t, script1.ID)
	if err != nil {
		t.Fail(err)
	}

	if script1.ID != script.ID {
		t.FailStr("create script and get script IDs don't match")
	}

	_, err = updateScript(t, script1.ID)
	if err != nil {
		t.Fail(err)
	}

	scripts, err := getScripts(t, script1.ID)
	if err != nil {
		t.Fail(err)
	}

	if len(scripts) < 2 {
		t.Fail(err)
	}

	err = deleteScript(t, script1.ID)
	if err != nil {
		t.Fail(err)
	}

	err = deleteScript(t, script2.ID)
	if err != nil {
		t.Fail(err)
	}

	return t
}

func createVersion(t *TestRun, id string) (*db.ScriptVersion, error) {

	body := "var x = 2"
	d, _ := json.Marshal(map[string]string{"body": body})
	route := fmt.Sprintf("/api/scripts/%s/versions", id)

	req := newRequest("POST", platformRoute(route), d)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var version db.ScriptVersion
	json.NewDecoder(res.Body).Decode(&version)

	if version.Body != body {
		return &version, errors.New("created version body do not match")
	}

	if res.StatusCode != http.StatusOK {
		return &version, ErrorStatusCodeNotOk
	}

	return &version, nil
}

func getVersions(t *TestRun, id string) ([]db.ScriptVersion, error) {

	route := fmt.Sprintf("/api/scripts/%s/versions", id)

	req := newRequest("GET", platformRoute(route), nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var versions []db.ScriptVersion
	json.NewDecoder(res.Body).Decode(&versions)

	if res.StatusCode != http.StatusOK {
		return versions, ErrorStatusCodeNotOk
	}

	return versions, nil
}

func CRUDVersion() *TestRun {

	t := &TestRun{Name: "CRUDVersion"}

	script, err := createScript(t)
	if err != nil {
		t.Fail(err)
	}

	_, err = createVersion(t, script.ID)
	if err != nil {
		t.Fail(err)
	}

	versions, err := getVersions(t, script.ID)
	if err != nil {
		t.Fail(err)
	}

	if len(versions) == 0 {
		t.FailStr("bad number of versions")
	}

	err = deleteScript(t, script.ID)
	if err != nil {
		t.Fail(err)
	}

	return t
}

func VersionsCountLimited() *TestRun {

	t := &TestRun{Name: "VersionsCountLimited"}

	script, err := createScript(t)
	if err != nil {
		t.Fail(err)
	}

	for i := 0; i < db.MaxScriptVersions+5; i++ {
		_, err = createVersion(t, script.ID)
		if err != nil {
			t.Fail(err)
		}
	}

	versions, err := getVersions(t, script.ID)
	if err != nil {
		t.Fail(err)
	}

	if len(versions) > db.MaxScriptVersions {
		t.FailStr("max script version count exceeded")
	}

	err = deleteScript(t, script.ID)
	if err != nil {
		t.Fail(err)
	}

	return t
}
