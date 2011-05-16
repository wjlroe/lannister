package main

import (
	"http"
	"fmt"
	"os"
	"strconv"
	"path/filepath"
	"io"
)

var content_types = map[string] string {
	"jpg" : "image/jpeg",
	"png" : "image/png",
	"gif" : "image/gif",
	"js" : "text/javascript",
	"css" : "text/css",
	"html" : "text/html",
}

func Serve() {
	// http.HandleFunc("/", FileRequest)
	http.Handle("/", http.FileServer("./", "/"))
	http.ListenAndServe(":6565", nil)
}

func FileRequest(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path
	fmt.Println("Filename requested: ", filename)
	filename_stat, err := os.Stat(filename)
	if err != nil {
		http.Error(w, err.String(), http.StatusNotFound)
		return
	}

	if filename_stat.IsDirectory() {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}
	defer fd.Close()

	content_type := content_types[filepath.Ext(filename)]
	w.Header().Set("Content-Type", content_type)
	w.Header().Set("Content-Length", strconv.Itoa64(filename_stat.Size))
	io.Copy(w, fd)
}