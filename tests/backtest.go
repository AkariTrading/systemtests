package tests

import (
	"fmt"
	"net/http"
	"time"

	"github.com/akaritrading/backtest/pkg/backtestclient"
	"github.com/akaritrading/libs/flag"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func (test Test) BackTest(t *TestRun) {

	body := ""

	req := map[string]interface{}{
		"body":     body,
		"exchange": "binance",
		"symbolA":  "BTC",
		"symbolB":  "TRY",
		"start":    time.Now().Add(-(time.Hour * 24 * 7)).Unix() * 1000,
		"end":      time.Now().Unix() * 1000,
		"balance":  map[string]float64{"BTC": 0, "TRY": 1000},
	}

	header := make(http.Header)
	header[sessionTokenHeader] = []string{sessiontoken}

	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("%s/ws", flag.PlatformHost()), header)
	if err != nil {
		t.Fail(errors.Wrap(err, "connecting to backtest route failed"))
		return
	}
	defer c.Close()

	if err := c.WriteJSON(req); err != nil {
		t.Fail(errors.Wrap(err, "writing to backtest route failed"))
	}

	var res backtestclient.BacktestResponse

	for {
		if err := c.ReadJSON(&res); err != nil {
			t.Fail(errors.Wrap(err, "reading from backtest failed"))
			return
		}
	}
}
