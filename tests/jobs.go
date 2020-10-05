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
		"symbolB":    "TRY",
		"type":       "cycledryrun",
		"balance":    map[string]interface{}{"BTC": 0.0, "TRY": 1000},
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

func trades() ([]db.Trade, error) {
	var trades []db.Trade
	return trades, NewRequest("GET", PlatformRoute("/api/trades"), nil, &trades)
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

		MarketBuy(BuyBalance())
		MarketSell(SellBalance())

		Print(Average(30))
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

	logs, err := getLogs(job.ID)
	if err != nil {
		t.Fail(errors.Wrap(err, "get logs failed"))
	}

	j, err := getJob(job.ID)
	if err != nil {
		t.Fail(err)
	}

	if len(logs) != 3 {
		t.FailStr("job did not produce three logs")
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

	var balance map[string]float64
	json.Unmarshal(util.StrToBytes(j.BalanceJSON), &balance)

	if balance["BTC"] != 0 {
		t.FailStr("TRY must be zero")
	}

	if balance["TRY"] < 997 {
		t.FailStr("TRY must be greater than approx 997")
	}

	trades, err := trades()
	if err != nil {
		t.Fail(errors.Wrap(err, "could not fetch trades"))
	}

	if len(trades) != 2 {
		t.FailStr(fmt.Sprintf("trades returned %d trades instead of 2", len(trades)))
	}

	err = stopJob(j.ID)
	if err != nil {
		t.Fail(errors.Wrap(err, "stopJobs failed"))
		return
	}
}
