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

const apiAddr = "http://127.0.0.1:8080"
const rootPath = "./uploads"

func selfDeleteFileHandler(store *httpd.Store) {
	cid := store.R.FormValue("cid")
	path := fsutil.ResolveBase(rootPath, cid)
	if err := os.Remove(path); err != nil {
		logger.Panic(err)
	}
	store.Respond200([]byte("ok"))
}

type fileInfo struct {
	Cid  string `json:"cid"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Time int64  `json:"time"`
}

func generateCid() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		logger.Panic(err)
	}
	return hex.EncodeToString(buf)
}

func createFilesHandler(store *httpd.Store) {
	reader, err := store.R.MultipartReader()
	if err != nil {
		logger.Panic(err)
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Panic(err)
		}
		if part.FormName() == "file" {
			var cid string
			for max := 5; max >= 0; max-- {
				cid = generateCid()
				if _, err := os.Stat(fsutil.ResolveBase(rootPath, cid)); os.IsNotExist(err) {
					break
				}
				if max == 0 {
					logger.Panic("cid generate failed")
				}
			}
			file, err := os.Create(fsutil.ResolveBase(rootPath, cid))
			if err != nil {
				logger.Panic(err)
			}
			defer file.Close()
			n, err := io.Copy(file, part)
			if err != nil {
				logger.Panic(err)
			}
			body, err := json.Marshal(fileInfo{cid, part.FileName(), n, time.Now().Unix()})
			if err != nil {
				logger.Panic(err)
			}
			_, err = http.Post(apiAddr+"/self/api/file", "application/json", bytes.NewBuffer(body))
			if err != nil {
				logger.Panic(err)
			}
			store.Redirect("/view/", http.StatusFound)
		}
	}
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
	mux := httpd.NewMux()
	mux.Handle("/self/file/data", "DELETE", selfDeleteFileHandler)
	mux.Handle("/file/data", "POST", createFilesHandler)
	mux.Handle("/file/data", "GET", downloadFileHandler)

	logger.Info("FILE service started.")
	if err := http.ListenAndServe(":8081", logger.Req(logger.Recovery(mux))); err != nil {
		logger.Fatal(err)
	}
}
