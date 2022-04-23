package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/glb/util/fsutil"
)

var (
	listenAddr = "0.0.0.0:8081"
	apiPrefix  = "http://127.0.0.1:8080"
	rootPath   = "./uploads"
)

func init() {
	if val, ok := os.LookupEnv("LISTEN_ADDR"); ok {
		listenAddr = val
	}
	if val, ok := os.LookupEnv("API_PREFIX"); ok {
		apiPrefix = val
	}
	if val, ok := os.LookupEnv("ROOT_PATH"); ok {
		rootPath = val
	}
}

type fileInfo struct {
	Cid  string
	Name string
	Size int64
	Time int64
}

func selfDeleteFileHandler(store *httpd.Store) {
	cid := store.R.FormValue("cid")
	path := fsutil.ResolveBase(rootPath, cid)
	if err := os.Remove(path); err != nil {
		logger.Panic(err)
	}
	store.Respond200([]byte("ok"))
}

func uploadFileHandler(store *httpd.Store) {
	// get file reader
	input, header, err := store.R.FormFile("file")
	if err != nil {
		logger.Panic(err)
	}
	defer input.Close()
	// generate cid
	buf := make([]byte, 16)
	_, err = rand.Read(buf)
	if err != nil {
		logger.Panic(err)
	}
	cid := hex.EncodeToString(buf)
	// check existing
	path := fsutil.ResolveBase(rootPath, cid)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		logger.Panic("random cid duplicated")
	}
	// save file
	file, err := os.Create(path)
	if err != nil {
		logger.Panic(err)
	}
	defer file.Close()
	n, err := io.Copy(file, input)
	if err != nil {
		logger.Panic(err)
	}
	// save meta data
	body, err := json.Marshal(fileInfo{cid, header.Filename, n, time.Now().Unix()})
	if err != nil {
		logger.Panic(err)
	}
	_, err = http.Post(apiPrefix+"/self/api/file", "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Panic(err)
	}
	store.Redirect("/view/", http.StatusFound)
}

func downloadFileHandler(store *httpd.Store) {
	cid := store.R.FormValue("cid")
	path := fsutil.ResolveBase(rootPath, cid)
	file, err := os.Open(path)
	if err != nil {
		logger.Panic(err)
	}
	defer file.Close()

	filename := store.R.FormValue("name")
	store.W.Header().Set("content-disposition", "attachment; filename*=UTF-8''"+filename+"; filename=\""+filename+"\"")

	http.ServeFile(store.W, store.R, path)
}

func main() {
	info, err := os.Stat(rootPath)
	if err == nil && !info.IsDir() {
		logger.Fatal("root path is not a directory")
	} else if os.IsNotExist(err) {
		logger.Info("create root directory")
		err = os.MkdirAll(rootPath, 0755)
	}
	if err != nil {
		logger.Fatal(err)
	}

	mux := httpd.NewMux()
	mux.Handle("/self/file/data", "DELETE", selfDeleteFileHandler)
	mux.Handle("/file/data", "POST", uploadFileHandler)
	mux.Handle("/file/data", "GET", downloadFileHandler)

	logger.Info("FILE service started.")
	if err := http.ListenAndServe(listenAddr, logger.Req(logger.Recovery(mux))); err != nil {
		logger.Fatal(err)
	}
}
