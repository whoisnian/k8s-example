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
	DatabaseURI string // database connection URI

	StorageDriver string // storage driver, filesystem or aws-s3
	RootDirectory string // filesystem: root directory
	S3Endpoint    string // aws-s3: endpoint
	S3AccessKey   string // aws-s3: access key id
	S3SecretKey   string // aws-s3: secret access key
	S3Secure      bool   // aws-s3: secure transport
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

		CFG.ListenAddr = stringFromEnv("CFG_LISTENADDR", "0.0.0.0:8081")
		CFG.DatabaseURI = stringFromEnv("CFG_DATABASEURI", "root:password@tcp(127.0.0.1:3306)/dbname")

		CFG.StorageDriver = stringFromEnv("CFG_STORAGEDRIVER", "filesystem")
		CFG.RootDirectory = stringFromEnv("CFG_ROOTDIRECTORY", "./uploads")
		CFG.S3Endpoint = stringFromEnv("CFG_S3ENDPOINT", "https://s3.amazonaws.com")
		CFG.S3AccessKey = stringFromEnv("CFG_S3ACCESSKEY", "QZH1XZPZLP8DA3GKA3J1")
		CFG.S3SecretKey = stringFromEnv("CFG_S3SECRETKEY", "VQyou21kIHVuKLkULNaETFnN7kLstyiX2KEtVbuI")
		CFG.S3Secure = boolFromEnv("CFG_S3SECURE", true)
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
