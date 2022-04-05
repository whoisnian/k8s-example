package main

import (
	"net/http"
	"sort"

	"github.com/whoisnian/glb/httpd"
	"github.com/whoisnian/glb/logger"
)

func selfCreateFileHandler(store *httpd.Store) {
	store.Respond200([]byte("ok"))
}

type fileInfo struct {
	Cid  string `json:"cid"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Time int64  `json:"time"`
}

type fileInfos []fileInfo

func (arr fileInfos) Len() int {
	return len(arr)
}
func (arr fileInfos) Less(i, j int) bool {
	return arr[i].Time > arr[j].Time
}
func (arr fileInfos) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

var mockMap = map[string]fileInfo{
	"aaaaaaaaaa": {"aaaaaaaaaa", "password.txt", 1743, 1649149438},
	"bbbbbbbbbb": {"bbbbbbbbbb", "run.exe", 17090020, 1649141438},
	"cccccccccc": {"cccccccccc", "image.jpg", 2939130, 1649149438},
	"dddddddddd": {"dddddddddd", "result.json", 88552, 1641949438},
	"eeeeeeeeee": {"eeeeeeeeee", "arch.iso", 906309632, 1645149438},
}

func listFilesHandler(store *httpd.Store) {
	var files fileInfos
	for _, v := range mockMap {
		files = append(files, v)
	}
	sort.Sort(files)
	store.RespondJson(files)
}

func deleteFileHandler(store *httpd.Store) {
	cid := store.R.FormValue("cid")
	delete(mockMap, cid)
	store.Respond200([]byte("ok"))
}

func main() {
	mux := httpd.NewMux()
	mux.Handle("/self/api/file", "POST", selfCreateFileHandler)
	mux.Handle("/api/files", "GET", listFilesHandler)
	mux.Handle("/api/file", "DELETE", deleteFileHandler)

	logger.Info("API service started.")
	if err := http.ListenAndServe(":8080", logger.Req(logger.Recovery(mux))); err != nil {
		logger.Fatal(err)
	}
}
