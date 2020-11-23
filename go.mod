module github.com/akaritrading/systemtests

go 1.15

require (
	github.com/akaritrading/backtest v0.0.0
	github.com/akaritrading/engine v0.0.0
	github.com/akaritrading/libs v0.0.0
	github.com/gorilla/websocket v1.4.2
	github.com/pkg/errors v0.9.1
)

replace github.com/akaritrading/backtest v0.0.0 => ../backtest

replace github.com/akaritrading/prices v0.0.0 => ../prices

replace github.com/akaritrading/engine v0.0.0 => ../engine

replace github.com/akaritrading/libs v0.0.0 => ../libs
