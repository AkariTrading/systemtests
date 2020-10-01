package tests

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akaritrading/libs/db"
)

func createVersion(id string) (*db.ScriptVersion, error) {

	body := "var x = 2"
	d, _ := json.Marshal(map[string]string{"body": "var x = 2"})

	var version db.ScriptVersion

	err := NewRequest("POST", PlatformRoute(fmt.Sprintf("/api/scripts/%s/versions", id)), d, &version)
	if err != nil {
		return nil, err
	}

	if version.Body != body {
		return &version, errors.New("created version body do not match")
	}

	return &version, nil
}

func getVersions(t *TestRun, id string) ([]db.ScriptVersion, error) {
	var versions []db.ScriptVersion
	return versions, NewRequest("GET", PlatformRoute(fmt.Sprintf("/api/scripts/%s/versions", id)), nil, &versions)
}

func (test Test) CRUDVersion(t *TestRun) {

	script, err := createScript()
	if err != nil {
		t.Fail(err)
	}

	_, err = createVersion(script.ID)
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

	err = deleteScript(script.ID)
	if err != nil {
		t.Fail(err)
	}
}

func (test Test) VersionsCountLimited(t *TestRun) {

	script, err := createScript()
	if err != nil {
		t.Fail(err)
	}

	for i := 0; i < db.MaxScriptVersions+5; i++ {
		_, err = createVersion(script.ID)
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

	err = deleteScript(script.ID)
	if err != nil {
		t.Fail(err)
	}
}
