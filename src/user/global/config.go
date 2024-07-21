package global

import (
	"encoding/json"
	"os"
)

var CFG Config

type Config struct {
	Debug   bool // enable debug output
	Version bool // show version and quit

	ListenAddr  string // server listen addr
	MysqlDSN    string // mysql dsn from https://github.com/go-sql-driver/mysql/blob/master/README.md#dsn-data-source-name
	AutoMigrate bool   // automatically migrate mysql schema and quit
	RedisURI    string // redis uri from https://github.com/redis/redis-specifications/blob/master/uri/redis.txt
	AppSecret   string // session authentication key with 32 bytes

	DisableRegistration bool // disable user self-registration

	TraceEndpointUrl string // OTLP Trace HTTP Exporter endpoint url
}

func SetupConfig() {
	if filename, ok := os.LookupEnv("CFG_CONFIG"); ok {
		fi, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer fi.Close()

		if err = json.NewDecoder(fi).Decode(&CFG); err != nil {
			panic(err)
		}
	} else {
		CFG.Debug = boolFromEnv("CFG_DEBUG", false)
		CFG.Version = boolFromEnv("CFG_VERSION", false)

		CFG.ListenAddr = stringFromEnv("CFG_LISTENADDR", "0.0.0.0:8080")
		CFG.MysqlDSN = stringFromEnv("CFG_MYSQLDSN", "root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=UTC")
		CFG.AutoMigrate = boolFromEnv("CFG_AUTOMIGRATE", false)
		CFG.RedisURI = stringFromEnv("CFG_REDISURI", "redis://default:password@127.0.0.1:6379/0")
		CFG.AppSecret = stringFromEnv("CFG_APPSECRET", "authentication_key_with_32_bytes")

		CFG.DisableRegistration = boolFromEnv("CFG_DISABLEREGISTRATION", false)

		CFG.TraceEndpointUrl = stringFromEnv("CFG_TRACEENDPOINTURL", "http://127.0.0.1:4318")
	}
}

func boolFromEnv(envKey string, defVal bool) bool {
	if val, ok := os.LookupEnv(envKey); !ok {
		return defVal
	} else if val == "" || val == "true" {
		return true
	} else if val == "false" {
		return false
	} else {
		panic("boolFromEnv: " + val)
	}
}

func stringFromEnv(envKey string, defVal string) string {
	if val, ok := os.LookupEnv(envKey); !ok {
		return defVal
	} else {
		return val
	}
}
