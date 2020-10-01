package tests

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akaritrading/libs/db"
)

func createExchange() (*db.ExchangeConnection, error) {

	d, _ := json.Marshal(map[string]string{
		"exchange":  "binance",
		"apiKey":    "UAbv4w94TT7E0vB2JfHS8gwxXQvUFDzBjlZ4vRxGZ6njrbOYX8oMH5QMW1BXJAA2",
		"apiSecret": "VsZwXQr68HwjSAI4kuZLHU2vsYsW6RjXztuuUwO25sP6Zyqt0N4bLRzoYQGyKU78",
	})

	var conn db.ExchangeConnection

	err := NewRequest("POST", PlatformRoute("/api/userExchanges/"), d, &conn)

	return &conn, err
}

func deleteExchange(exchangeID string) error {
	return NewRequest("DELETE", PlatformRoute(fmt.Sprintf("/api/userExchanges/%s", exchangeID)), nil, nil)
}

func getExchanges(id string) ([]db.ExchangeConnection, error) {

	var conn []db.ExchangeConnection

	err := NewRequest("GET", PlatformRoute("/api/userExchanges"), nil, &conn)

	if err != nil {
		return nil, err
	}

	for _, c := range conn {
		if c.ID == id {
			return conn, nil
		}
	}

	return conn, errors.New("could not find exchange with ID in returned list")
}

func (test Test) ConnectExchange(t *TestRun) {

	exc, err := createExchange()
	if err != nil {
		t.Fail(err)
	}

	if exc.Exchange != "binance" {
		t.FailStr("exchange name does not match")
	}

	if exc.UserID != user.ID {
		t.FailStr("exchange user ID does not match")
	}

	exchages, err := getExchanges(exc.ID)
	if err != nil {
		t.Fail(err)
	}

	if len(exchages) < 1 {
		t.FailStr("fetching connecting exchanges returned empty list")
	}

	err = deleteExchange(exc.ID)
	if err != nil {
		t.Fail(err)
	}
}
