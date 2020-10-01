package tests

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akaritrading/engine/pkg/engineclient"
	"github.com/akaritrading/libs/db"
	"github.com/akaritrading/libs/util"
	"github.com/pkg/errors"
)

func runDryRun(scriptID string, scriptBody string, exchangeID string) (*engineclient.JobRequest, error) {

	body := map[string]interface{}{
		"body":       scriptBody,
		"exchange":   "binance",
		"exchangeId": exchangeID,
		"symbolA":    "BTC",
		"symbolB":    "USDT",
		"type":       "cycledryrun",
		"balance":    map[string]interface{}{"BTC": 0.0, "USDT": 1000},
		"scriptId":   scriptID,
	}
	d, _ := json.Marshal(body)

	var res engineclient.JobRequest
	return &res, NewRequest("POST", PlatformRoute("/api/jobs"), d, &res)
}

func getLogs(jobID string) ([]db.JobLog, error) {
	var logs []db.JobLog
	return logs, NewRequest("GET", PlatformRoute(fmt.Sprintf("/api/jobs/%s/logs", jobID)), nil, &logs)
}

func getJob(jobID string) (*db.Job, error) {
	var job db.Job
	return &job, NewRequest("GET", PlatformRoute(fmt.Sprintf("/api/jobs/%s", jobID)), nil, &job)
}

func stopJob(jobID string) error {
	return NewRequest("DELETE", PlatformRoute(fmt.Sprintf("/api/jobs/%s", jobID)), nil, nil)
}

func (test Test) Jobs(t *TestRun) {

	script, err := createScript()
	if err != nil {
		t.Fail(err)
	}
	defer deleteScript(script.ID)

	ex, err := createExchange()
	if err != nil {
		t.Fail(err)
		return
	}
	defer deleteExchange(ex.ID)

	body := `
	
		var s = GetState();
		s.test = 1;
		SaveState(s)

		MarketBuy(Balance().USDT)
		MarketSell(Balance().BTC)

		Print(OrderbookPrice().sell)		
		Print(OrderbookPrice().buy)		
	`
	job, err := runDryRun(script.ID, body, ex.ID)
	if err != nil {
		t.Fail(errors.New("dry run failed"))
		return
	}

	if job.ID == "" {
		t.FailStr("dry run returned empty ID")
		return
	}

	time.Sleep(time.Second * 2)

	l, err := getLogs(job.ID)
	if err != nil {
		t.Fail(errors.Wrap(err, "get logs failed"))
		return
	}

	j, err := getJob(job.ID)
	if err != nil {
		t.Fail(err)
		return
	}

	if len(l) != 2 {
		t.FailStr("job did not produce two logs")
		return
	}

	if j.Body != body {
		t.FailStr("missing fields in created job")
	}

	var state map[string]interface{}
	json.Unmarshal(util.StrToBytes(j.StateJSON), &state)

	if state["test"].(float64) != 1 {
		t.FailStr("state was not saved properly")
	}

	if len(j.Trades) != 2 {
		t.FailStr("job did not produce two trades")
	}

	err = stopJob(j.ID)
	if err != nil {
		t.Fail(errors.Wrap(err, "stopJobs failed"))
		return
	}
}
