package main

import (
	"github.com/russross/blackfriday"
	"os"
	"io/ioutil"
	"io"
	"fmt"
	"path"
	"net/http"
	"path/filepath"
	htmpl "html/template"
	"strings"
	"log"
)

// TODO: Add extra URIs in case of failure
//http://code.jquery.com/jquery-1.5.2.min.js
// TODO: Make generic func url -> location downloader

var static_files = map[string] string {
	"http://ajax.googleapis.com/ajax/libs/jquery/1.5.2/jquery.min.js" : "javascript",
	"http://ajax.googleapis.com/ajax/libs/jqueryui/1.8.12/jquery-ui.min.js" : "javascript",
	"https://github.com/defunkt/jquery-pjax/raw/master/jquery.pjax.js" : "javascript",
	"http://d3nwyuy0nl342s.cloudfront.net/images/modules/facebox/loading.gif" : "images",
}

type Page struct {
	PageContent string
}

func download(url string, location string) {
	filename := path.Base(url)
	fullname := path.Join(location, filename)
	os.MkdirAll(location, 0755)
	//fmt.Printf("Fullname: %s\n", fullname)
	encoded := fmt.Sprintf("%s", url)
	///fmt.Println(encoded)
	r, err := http.Get(encoded)
	const NBUF = 512
	var buf [NBUF]byte
	fd, oerr := os.OpenFile(fullname, os.O_WRONLY|os.O_CREATE, 0644)
	if oerr != nil {
		fmt.Printf("Opening file: %s failed with error: %s\n", fullname, oerr.Error())
		os.Exit(1)
	}
	defer fd.Close()
	if err == nil {
		defer r.Body.Close()
		for {
			nr, ferr := r.Body.Read(buf[0:])
			if ferr != nil {
				if ferr == io.EOF {
					//fmt.Println("EOF")
					break
				} else {
					fmt.Printf("Error: %s  downloading uri: %s\n", ferr, encoded)
					os.Exit(1)
				}
			}
			nw, ew := fd.Write(buf[0:nr])
			if ew != nil {
				fmt.Printf("Error writing to file. Error: %s\n", ew)
				os.Exit(1)
			}
			if nw != nr {
				fmt.Printf("Error writing %d bytes from Download!\n", nr)
				os.Exit(1)
			}
		}
		//fmt.Printf("Finished reading/writing file\n")

	} else {
		fmt.Printf("Error reading from body: %s\n", err)
		//log.Stderr(err)
		os.Exit(1)
	}
	//fmt.Println("Written. Closing...")

	//fmt.Println("Closed")
}

func is_dir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func checkDirectory(root string) {
	dirs := []string{"pages","layouts","site"}
	for _, dir := range dirs {
		if !is_dir(filepath.Join(root, dir)) {
			fmt.Println("Current directory is not a Lannister site dir.")
			fmt.Println("To be used, there should be ./pages/, ./layouts/, ./site/")
			fmt.Println("Please use ./lannister createsite ./dir/to/create/as/site")
			os.Exit(1)
		}
	}
}

// TODO: Return the markdown doc so it can be used multiple times instead of parsing it again
func markdown_parse(file_in io.Reader) []byte {
	b, _ := ioutil.ReadAll(file_in)

	return blackfriday.MarkdownCommon(b)
}

func write_file(content string, filepath string) {
	out_fd, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Failed to open file: %s for writing, error: %s\n", filepath, err)
		os.Exit(1)
	}
	defer out_fd.Close()
	out_fd.WriteString(content)
}

func createsite(site_dir string) {
	os.MkdirAll(filepath.Join(site_dir, "pages"), 0755)
	os.MkdirAll(filepath.Join(site_dir, "site"), 0755)
	os.MkdirAll(filepath.Join(site_dir, "layouts"), 0755)
	for uri, path := range static_files {
		local_path := filepath.Join(site_dir, path)
		fmt.Printf("Downloading URI: %s to path: %s\n", uri, local_path)
		download(uri, local_path)
	}
	appjs_path := filepath.Join(site_dir, "javascript", "app.js")
	write_file(appjs, appjs_path)
	index_path := filepath.Join(site_dir, "pages", "index.md")
	write_file(index_page, index_path)
	about_path := filepath.Join(site_dir, "pages", "about.md")
	write_file(about_page, about_path)
	default_path := filepath.Join(site_dir, "layouts", "default.html")
	write_file(layout_default, default_path)
	pjax_path := filepath.Join(site_dir, "layouts", "default-pjax.html")
	write_file(layout_pjax, pjax_path)
	// TODO: default.rss ?
}

func CopyFile(dst, src string) (int64, error) {
        sf, err := os.Open(src)
        if err != nil {
                return 0, err
        }
        defer sf.Close()
        df, err := os.Create(dst)
        if err != nil {
                return 0, err
        }
        defer df.Close()
        return io.Copy(df, sf)
}

func copy_dir_contents(root, src_dir string) {
	dst_dir := filepath.Join(root, "site", src_dir)
	src_fd, err := os.Open(filepath.Join(root, src_dir))
	if err != nil {
		log.Fatalf("Could not open dir: %s", src_dir)
	}
	files, err := src_fd.Readdirnames(-1)
	if err != nil {
		log.Fatalf("Could not read directory names: %s", err)
	}
	for _,filename := range files {
		dst_file := filepath.Join(dst_dir, filename)
		src_file := filepath.Join(root, src_dir, filename)
		log.Printf("Copying %s to %s", src_file, dst_file)
		_, err := CopyFile(dst_file, src_file)
		if err != nil {
			log.Fatalf("Error: %s writing to %s!", err, dst_file)
		}
	}
}

func generate(root string) error {
	// list the files to be templated
	// process them - markdown
	// run each through each layout file (TODO: make configurable)
	// save each output file in the site directory
	os.MkdirAll(filepath.Join(root, "site", "images"), 0755)
	os.MkdirAll(filepath.Join(root, "site", "javascript"), 0755)
	copy_dir_contents(root, "images")
	copy_dir_contents(root, "javascript")

	var page_files []string
	var err error
	page_files, err = filepath.Glob(filepath.Join(root, "pages", "*.md"))
	if err != nil {
		fmt.Printf("Failed to find any .md files in pages subdir. Error: %s\n", err)
		return err
	}

	var template_files []string
	template_files, err = filepath.Glob(filepath.Join(root, "layouts", "*.html"))
	if err != nil {
		fmt.Printf("Failed to find any applicable layout files. Error: %s\n", err)
		return err
	}
	// var fmap = template.FormatterMap {
	// 	"" : template.StringFormatter,
	// 	"html": template.HTMLFormatter,
	// }
	templates := map[string] *htmpl.Template{}
	for _,t := range template_files {
		templates[filepath.Base(t)] = htmpl.Must(htmpl.ParseFiles(t))
	}

	for _,page_filepath := range page_files {
		in_fd, err := os.Open(page_filepath)
		if err != nil {
			fmt.Printf("Failed to open file: %s, error: %s\n", page_filepath, err)
			return err
		}
		defer in_fd.Close()
		doc := markdown_parse(in_fd)
		page := &Page{PageContent: string(doc)}
		page_filename := filepath.Base(page_filepath)
		page_ext := filepath.Ext(page_filename)
		fmt.Printf("Page ext: %s\n", page_ext)
		page_filename = strings.Replace(page_filename, page_ext, "", -1)
		for tpl_filename, template := range templates {
			// TODO: This should be pagename-pjax.html or pagename.html - BUG
			filename := strings.Replace(tpl_filename, "default", page_filename, -1)
			fmt.Printf("tpl_filename: %s, page_filename: %s, Output page filename: %s\n", tpl_filename, page_filename, filename)
			filepath := filepath.Join(root, "site", filename)
			fmt.Printf("Going to save templated page: %s as file: %s\n", page_filename, filepath)
			out_fd, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				fmt.Printf("Failed to open file: %s for writing, error: %s\n", filepath, err)
				return err
			}
			defer out_fd.Close()
			template.Execute(out_fd, page)
			out_fd.Sync()
		}
	}
	return nil
}

func GetCWD() (directory string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("OMG can't work out the CWD")
	}
	return dir
}

func main() {
	// Get the command
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "createsite":
			fmt.Println("Createsite")
			if len(os.Args) > 2 {
				directory := os.Args[2]
				fmt.Printf("Directory: %s\n", directory)
				createsite(directory)
			} else {
				fmt.Println("Please provide a directory to create the site in")
			}
		case "generate":
			directory := ""
			if len(os.Args) > 2 {
				directory = os.Args[2]
			} else {
				directory = GetCWD()
			}
			checkDirectory(directory)
			generate(directory)
		case "serve":
			if len(os.Args) > 2 {
				if os.Args[2] == "dropbox" {
					DropBoxServe()
				} else {
					fmt.Println("Don't understand what you want to serve!. Did you mean 'dropbox'?")
				}
			} else {
				Serve("./")
			}
		default:
			fmt.Printf("Command: %s not understood.\n", command)
		}
	} else {
		fmt.Println("no command specified!")
	}
}

const appjs = `
$(document).ready(function() {
		      $('nav a').pjax('#main', {
						beforeSend: function(xhr){
						    xhr.setRequestHeader('X-PJAX', 'true');
						    this.url = this.url.replace("^/$", "/index-pjax.html");
						    this.url = this.url.replace(".html", "-pjax.html");
						}
					    });

		      $('#main')
			  .bind('start.pjax', function() {
				    console.log("start pjax");
				    $('#main').hide("slide", {direction: "left"}, 1000);
				    $('#loading').show();
				})
			  .bind('pjax', function() {
				console.log("pjax fired");
				})
			  .bind('end.pjax', function() {
				    $('#loading').hide();
				    $('#main').show("slide", {direction: "right"}, 1000);
				});
		  });
`

const about_page = `<article>
## About
This is the about page.

</article>`

const index_page = `<article>
## Index
This is the index page.

</article>`

const layout_pjax = `
{PageContent}
`

const layout_default = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <script src="/javascript/jquery.min.js"></script>
    <script src="/javascript/jquery-ui.min.js"></script>
    <script src="/javascript/jquery.pjax.js"></script>
    <script src="/javascript/app.js"></script>
  </head>
  <body>
    <header>
      <h1>Site title</h1>
      <nav>
        <a href="/index.html">Home</a>
        <a href="/about.html">About</a>
      </nav>
    </header>
    <div id="loading" style="display:none">
      <img src="/images/loading.gif" />
    </div>
    <div id="main">
      {PageContent}
    </div>
  </body>
</html>
`
