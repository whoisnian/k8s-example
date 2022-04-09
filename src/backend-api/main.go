package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
)

var (
	listenAddr = "0.0.0.0:8080"
	filePrefix = "http://127.0.0.1:8081"
	mysqlDSN   = "k8s:KxY8cSAWz1WJEfs3@tcp(127.0.0.1:3306)/k8s"

	db *sql.DB
)

func init() {
	if val, ok := os.LookupEnv("LISTEN_ADDR"); ok {
		listenAddr = val
	}
	if val, ok := os.LookupEnv("FILE_PREFIX"); ok {
		filePrefix = val
	}
	if val, ok := os.LookupEnv("MYSQL_DSN"); ok {
		mysqlDSN = val
	}
}

type fileInfo struct {
	Cid  string
	Name string
	Size int64
	Time int64
}

const createTableSQL = `
CREATE TABLE IF NOT EXISTS files (
    cid VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    size INT NOT NULL,
    time INT NOT NULL,
    PRIMARY KEY (cid)
)`

func selfCreateFileHandler(store *httpd.Store) {
	var info fileInfo
	err := json.NewDecoder(store.R.Body).Decode(&info)
	if err != nil {
		logger.Panic(err)
	}
	_, err = db.Exec("INSERT files SET cid=?, name=?, size=?, time=?", info.Cid, info.Name, info.Size, info.Time)
	if err != nil {
		logger.Panic(err)
	}
	store.Respond200([]byte("ok"))
}

func listFilesHandler(store *httpd.Store) {
	var files []fileInfo
	rows, err := db.Query("SELECT cid,name,size,time FROM files ORDER BY time DESC")
	if err != nil {
		logger.Panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var info fileInfo
		if err := rows.Scan(&info.Cid, &info.Name, &info.Size, &info.Time); err != nil {
			logger.Panic(err)
		}
		files = append(files, info)
	}
	store.RespondJson(files)
}

func deleteFileHandler(store *httpd.Store) {
	cid := store.R.FormValue("cid")
	_, err := db.Exec("DELETE FROM files WHERE cid=?", cid)
	if err != nil {
		logger.Panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", filePrefix+"/self/file/data?cid="+url.QueryEscape(cid), nil)
	if err != nil {
		logger.Panic(err)
	}
	_, err = client.Do(req)
	if err != nil {
		logger.Panic(err)
	}
	store.Respond200([]byte("ok"))
}

func main() {
	var err error
	db, err = sql.Open("mysql", mysqlDSN)
	if err != nil {
		logger.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()
	db.Exec(createTableSQL)

	mux := httpd.NewMux()
	mux.Handle("/self/api/file", "POST", selfCreateFileHandler)
	mux.Handle("/api/files", "GET", listFilesHandler)
	mux.Handle("/api/file", "DELETE", deleteFileHandler)

	logger.Info("API service started.")
	if err := http.ListenAndServe(listenAddr, logger.Req(logger.Recovery(mux))); err != nil {
		logger.Fatal(err)
	}
}
