package main

import (
	"net/http"
	"fmt"
	"os"
	"strconv"
	"path/filepath"
	"io"
	"log"
	"regexp"
	"strings"
)

var content_types = map[string] string {
	"jpg" : "image/jpeg",
	"png" : "image/png",
	"gif" : "image/gif",
	"js" : "text/javascript",
	"css" : "text/css",
	"html" : "text/html",
}

var href_regex = regexp.MustCompile("(href=\"/[^\"]*\")")

func Serve(root string) {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(root))))
	if err := http.ListenAndServe(":6565", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
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
		if site_info.IsDir() {
			site_prefix := "/" + sitename + "/"
			server := &DropboxServer{location: root, site_prefix: sitename, dropbox_location: dropbox_root}
			fmt.Printf("Serving %s at http://localhost:6767/%s/ - site_prefix: %s\n", sitename, sitename, site_prefix)
			//http.Handle(site_prefix, http.FileServer(root, site_prefix))
			http.Handle(site_prefix, server)
			num_dirs++
		}
	}
	if num_dirs > 0 {
		fmt.Println("Listening on http://localhost:6767")
		http.HandleFunc("/", LogRequest)
		http.ListenAndServe(":6767", nil)
	} else {
		log.Fatal("No directories existed to serve")
	}
}

type DropboxServer struct {
	location string
	site_prefix string
	dropbox_location string
}

func (server *DropboxServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Attempting dropbox file")
	filename := r.URL.Path
	full_filename := server.dropbox_location + filename
	fmt.Println("Filename requested: ", full_filename)
	//fmt.Fprintf(w, "hi\n")
	FileRequest(full_filename, w, r)
}

func LogRequest(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Path
	fmt.Println("Default handler: ", filename)
	fmt.Fprintf(w, "hi\n")
}

func FileRequest(filename string, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Filename requested: ", filename)
	filename_stat, err := os.Stat(filename)
	if err != nil {
		fmt.Println("404 Not Found: ", filename)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if filename_stat.IsDir() {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fd.Close()

	content_type := content_types[filepath.Ext(filename)]

	if strings.Contains(content_type, "text") {

	}

	w.Header().Set("Content-Type", content_type)
	w.Header().Set("Content-Length", strconv.FormatInt(filename_stat.Size(), 10))
	io.Copy(w, fd)
}
