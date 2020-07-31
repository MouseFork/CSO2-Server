module github.com/KouKouChan/CSO2-Server

go 1.14

require (
	github.com/garyburd/redigo v1.6.0
	github.com/mattn/go-sqlite3 v1.14.0
	golang.org/x/text v0.3.3
	gopkg.in/ini.v1 v1.57.0
)

replace golang.org/x/text v0.3.3 => github.com/golang/text v0.3.3
