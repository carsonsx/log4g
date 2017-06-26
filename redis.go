package log4g

import (
	"encoding/json"
	"github.com/go-redis/redis"
)

func newRedisLogger(prefix string, flag int, lc *loggerConfig) Logger {
	redisLogger := new(RedisLogger)
	redisLogger.cli = redis.NewClient(&redis.Options{
		Addr:     lc.Address,
		Password: lc.Password,
		DB:       lc.DB,
	})
	redisLogger.lc = lc
	redisLogger.GenericLogger = newLogger(prefix, flag, redisLogger)
	return redisLogger
}

type RedisLogger struct {
	*GenericLogger
	cli *redis.Client
	lc  *loggerConfig
}

func (l *RedisLogger) Write(p []byte) (n int, err error) {

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

	if l.lc.RedisType == "list" {
		cmd := l.cli.RPush(l.lc.RedisKey, p)
		return int(cmd.Val()), cmd.Err()
	}
	return 0, nil
}

func (l *RedisLogger) Close() {
	l.cli.Close()
}
