package main

import (
	"http"
	"fmt"
	"os"
	"strconv"
	"path/filepath"
	"io"
	"log"
)

var content_types = map[string] string {
	"jpg" : "image/jpeg",
	"png" : "image/png",
	"gif" : "image/gif",
	"js" : "text/javascript",
	"css" : "text/css",
	"html" : "text/html",
}

func Serve(root string) {
	// http.HandleFunc("/", FileRequest)
	http.Handle("/", http.FileServer(root, "/"))
	http.ListenAndServe(":6565", nil)
}

func DropBoxServe() {
	// TODO: Use actual HOME env
	dropbox_root := "/Users/will/Dropbox/Sites"
	fmt.Println("Going to serve directories in Dropbox/Sites")
	dropbox_fd, err := os.Open(dropbox_root)
	if err != nil {
		log.Fatal("Cannot read your Dropbox/Sites folder - does it exist?")
	}
	files, err := dropbox_fd.Readdirnames(-1)
	if err != nil {
		log.Fatal("cannot list files in your Dropbox/Sites folder")
	}
	num_dirs := 0
	for _,sitename := range files {
		root := filepath.Join(dropbox_root, sitename)
		site_info, err := os.Stat(root)
		if err != nil {
			fmt.Printf("Can't stat dir: %s\n", root)
		}
		if site_info.IsDirectory() {
			site_prefix := "/" + sitename
			fmt.Printf("Serving %s at http://localhost:6767/%s - site_prefix: %s\n", sitename, sitename, site_prefix)
			http.Handle(site_prefix, http.FileServer(root, site_prefix))
			num_dirs++
		}
	}
	if num_dirs > 0 {
		fmt.Println("Listening on http://localhost:6767")
		http.ListenAndServe(":6767", nil)
	} else {
		log.Fatal("No directories existed to serve")
	}
}

func ServeFiles(

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