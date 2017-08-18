package log4g

import (
	"encoding/json"
	"fmt"
	"net"
)

func newSocketLogger(level Level, prefix string, flag int, lc *loggerConfig, calldepth int) Logger {
	socketLogger := new(SocketLogger)
	var err error
	if lc.Network == "" {
		lc.Network = "udp"
	}
	socketLogger.conn, err = net.Dial(lc.Network, lc.Address)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	socketLogger.lc = lc
	socketLogger.GenericLogger = newLogger(level, prefix, flag, socketLogger, calldepth)
	return socketLogger
}

type SocketLogger struct {
	*GenericLogger
	conn net.Conn
	lc   *loggerConfig
}

func (l *SocketLogger) Write(p []byte) (n int, err error) {

	if l.conn == nil {
		return
	}

	if p[len(p)-1] == '\n' {
		p = p[0 : len(p)-1]
	}

	if l.lc.Codec == "json" {
		rec := make(map[string]interface{})
		rec[l.lc.JsonKey] = string(p)
		if l.lc.JsonExt != "" {
			var kv map[string]interface{}
			json.Unmarshal([]byte(l.lc.JsonExt), &kv)
			for k, v := range kv {
				rec[k] = v
			}
		}
		p, _ = json.Marshal(rec)
	}

	if l.lc.Network == "tcp" {
		p = append(p, '\n')
	}

	return l.conn.Write(p)
}

func (l *SocketLogger) Close() {
	if l.conn != nil {
		l.conn.Close()
	}
}
