package global

import (
	"os"
)

var (
	CFG_DEBUG   bool // enable debug output
	CFG_VERSION bool // show version and quit

	CFG_LISTENADDR  string // server listen addr
	CFG_DATABASEURI string // database connection URI

	CFG_STORAGE      string // storage mode, file or s3
	CFG_FILE_ROOT    string // file-storage: root directory
	CFG_S3_ENDPOINT  string // s3-storage: endpoint
	CFG_S3_ACCESSKEY string // s3-storage: access key id
	CFG_S3_SECRETKEY string // s3-storage: secret access key
	CFG_S3_SECURE    bool   // s3-storage: secure transport
)

func SetupConfig() {
	CFG_DEBUG = boolFromEnv("CFG_DEBUG", false)
	CFG_VERSION = boolFromEnv("CFG_VERSION", false)

	CFG_LISTENADDR = stringFromEnv("CFG_LISTENADDR", "0.0.0.0:8081")
	CFG_DATABASEURI = stringFromEnv("CFG_DATABASEURI", "root:password@tcp(127.0.0.1:3306)/dbname")

	CFG_STORAGE = stringFromEnv("CFG_STORAGE", "file")
	CFG_FILE_ROOT = stringFromEnv("CFG_FILE_ROOT", "./uploads")
	CFG_S3_ENDPOINT = stringFromEnv("CFG_S3_ENDPOINT", "https://s3.amazonaws.com")
	CFG_S3_ACCESSKEY = stringFromEnv("CFG_S3_ACCESSKEY", "BKJIK5AABMIBU2OMRH6B")
	CFG_S3_SECRETKEY = stringFromEnv("CFG_S3_SECRETKEY", "V7f1wqCQAc0woUEc8IEj5gUJVS5SQxhQo9SGr1r2")
	CFG_S3_SECURE = boolFromEnv("CFG_S3_SECURE", true)
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
