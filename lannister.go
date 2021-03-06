package main

import (
	"encoding/xml"
	"fmt"
	"github.com/russross/blackfriday"
	"golang.org/x/tools/blog/atom"
	"gopkg.in/yaml.v2"
	htmpl "html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type postMeta map[interface{}]interface{}

// Post object
type Post struct {
	Filename   string
	SourcePath string
	Title      string
	Metadata   postMeta
}

type page struct {
	PageContent htmpl.HTML
}

type destination interface {
	DestPath() string
}

// Date returns the publication date for a post
func (p *Post) Date() (t time.Time, err error) {
	const dateFormat = "2006-01-02 15:04"
	dateStr := p.Metadata["date"].(string)
	t, err = time.Parse(dateFormat, dateStr)
	if err != nil {
		fmt.Printf("Error parsing metadata date: %s. Error: %s\n", dateStr, err)
		return
	}
	return
}

// DestPath returns the translated filepath as it should be available in the site
func (p *Post) DestPath() (dest string, err error) {
	const datePathFormat = "2006/01/02"
	var dateRe = regexp.MustCompile(`^[\d\-]*`)
	filename := dateRe.ReplaceAllString(filepath.Base(p.Filename), "")

	var date time.Time
	date, err = p.Date()
	if err != nil {
		return
	}
	dest = filepath.Join(date.Format(datePathFormat), filename)

	return
}

func isDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func checkDirectory(root string) {
	dirs := []string{"pages", "layouts", "site"}
	for _, dir := range dirs {
		if !isDir(filepath.Join(root, dir)) {
			fmt.Println("Current directory is not a Lannister site dir.")
			fmt.Println("To be used, there should be ./pages/, ./layouts/, ./site/")
			fmt.Println("Please use ./lannister createsite ./dir/to/create/as/site")
			os.Exit(1)
		}
	}
}

// TODO: Return the markdown doc so it can be used multiple times instead of parsing it again
func markdownParse(fileIn io.Reader) []byte {
	b, _ := ioutil.ReadAll(fileIn)

	return blackfriday.MarkdownCommon(b)
}

func writeFile(content string, filepath string) {
	outFd, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Failed to open file: %s for writing, error: %s\n", filepath, err)
		os.Exit(1)
	}
	defer outFd.Close()
	outFd.WriteString(content)
}

func createAtomFeed(filename string, posts []*Post) (err error) {
	feed := atom.Feed{
		Title: "Will Roe's blog",
	}
	for _, p := range posts {
		var t time.Time
		t, err = p.Date()
		if err != nil {
			fmt.Printf("Error getting date of post: %v\n", p)
			break
		}
		e := &atom.Entry{
			Title:     p.Title,
			Published: atom.Time(t),
		}
		feed.Entry = append(feed.Entry, e)
	}

	var data []byte
	data, err = xml.Marshal(&feed)
	if err != nil {
		fmt.Printf("Failed to marshal the feed: %s", err)
		os.Exit(1)
	}
	writeFile(string(data[:]), filename)
	return
}

func createsite(siteDir string) {
	os.MkdirAll(filepath.Join(siteDir, "pages"), 0755)
	os.MkdirAll(filepath.Join(siteDir, "site"), 0755)
	os.MkdirAll(filepath.Join(siteDir, "layouts"), 0755)
	os.MkdirAll(filepath.Join(siteDir, "posts"), 0755)
	os.MkdirAll(filepath.Join(siteDir, "javascript"), 0755)
	os.MkdirAll(filepath.Join(siteDir, "images"), 0755)
	appjsPath := filepath.Join(siteDir, "javascript", "app.js")
	writeFile(appjs, appjsPath)
	indexPath := filepath.Join(siteDir, "pages", "index.md")
	writeFile(indexPage, indexPath)
	aboutPath := filepath.Join(siteDir, "pages", "about.md")
	writeFile(aboutPage, aboutPath)
	defaultPath := filepath.Join(siteDir, "layouts", "default.html")
	writeFile(layoutDefault, defaultPath)
	pjaxPath := filepath.Join(siteDir, "layouts", "default-pjax.html")
	writeFile(layoutPjax, pjaxPath)
	examplePostPath := filepath.Join(siteDir, "posts", "2015-02-27-example-post.md")
	writeFile(examplePost, examplePostPath)
	exampleConfigPath := filepath.Join(siteDir, "config.yml")
	writeFile(exampleConfig, exampleConfigPath)
}

func copyFile(dst, src string) (int64, error) {
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

func copyDirContents(root, srcDir string) {
	dstDir := filepath.Join(root, "site", srcDir)
	srcFd, err := os.Open(filepath.Join(root, srcDir))
	if err != nil {
		log.Fatalf("Could not open dir: %s", srcDir)
	}
	files, err := srcFd.Readdirnames(-1)
	if err != nil {
		log.Fatalf("Could not read directory names: %s", err)
	}
	for _, filename := range files {
		dstFile := filepath.Join(dstDir, filename)
		srcFile := filepath.Join(root, srcDir, filename)
		log.Printf("Copying %s to %s", srcFile, dstFile)
		_, err := copyFile(dstFile, srcFile)
		if err != nil {
			log.Fatalf("Error: %s writing to %s!", err, dstFile)
		}
	}
}

func getTemplates(root string) (templates map[string]*htmpl.Template, err error) {
	var templateFiles []string
	templateFiles, err = filepath.Glob(filepath.Join(root, "layouts", "*.html"))
	if err != nil {
		fmt.Printf("Failed to find any applicable layout files. Error: %s\n", err)
		return
	}
	templates = map[string]*htmpl.Template{}
	for _, t := range templateFiles {
		templates[filepath.Base(t)] = htmpl.Must(htmpl.ParseFiles(t))
	}
	return
}

func templatedPage(root string, sourcePath string, destinationPath string) error {
	templates, err := getTemplates(root)
	if err != nil {
		fmt.Printf("Failed to find templates: %s\n", err)
	}

	inFd, err := os.Open(filepath.Join(root, sourcePath))
	if err != nil {
		fmt.Printf("Failed to open file: %s, error: %s\n", sourcePath, err)
		return err
	}
	defer inFd.Close()

	fmt.Printf("\n\nRoot: %s Page: %s Dest page: %s\n", root, sourcePath, destinationPath)
	doc := markdownParse(inFd)
	page := &page{PageContent: htmpl.HTML(doc)}
	pageExt := filepath.Ext(destinationPath)
	pageFilename := strings.Replace(filepath.Base(destinationPath), pageExt, "", -1)
	for tplFilename, template := range templates {
		// TODO: This should be pagename-pjax.html or pagename.html - BUG
		filename := strings.Replace(tplFilename, "default", pageFilename, -1)
		fmt.Printf("tpl_filename: %s, page_filename: %s, Output page filename: %s\n", tplFilename, pageFilename, filename)
		filename = filepath.Join(root, "site", filepath.Dir(destinationPath), filename)
		os.MkdirAll(filepath.Dir(filename), 0755)
		fmt.Printf("Going to save templated page: %s as file: %s\n", pageFilename, filename)
		outFd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Printf("Failed to open file: %s for writing, error: %s\n", filename, err)
			return err
		}
		defer outFd.Close()
		template.Execute(outFd, page)
		outFd.Sync()
	}
	return nil
}

func templatedPost(root string, p *Post) error {
	dest, err := p.DestPath()
	if err != nil {
		fmt.Printf("Error with DestPath: %s\n", err)
		return err
	}
	return templatedPage(root, p.SourcePath, dest)
}

func relativePaths(elems ...string) (newPaths []string, err error) {
	var paths []string
	paths, err = filepath.Glob(filepath.Join(elems...))
	var newPath string
	for _, path := range paths {
		newPath, err = filepath.Rel(elems[0], path)
		if err != nil {
			return
		}
		newPaths = append(newPaths, newPath)
	}
	return
}

func postsMetadata(root string, postFiles []string) (posts []*Post, err error) {
	var data []byte
	var meta postMeta
	for _, postFilename := range postFiles {
		data, err = ioutil.ReadFile(filepath.Join(root, postFilename))
		if err != nil {
			fmt.Printf("Failed to open %s. Error: %s\n", postFilename, err)
			return
		}

		err = yaml.Unmarshal(data, &meta)
		if err != nil {
			fmt.Printf("Failed to unmarshal yaml from %s.Error: %s\n", postFilename, err)
		}

		p := new(Post)
		p.Metadata = meta
		p.Title = meta["title"].(string)
		p.Filename = filepath.Base(postFilename)
		p.SourcePath = postFilename
		posts = append(posts, p)
	}
	return
}

func generate(root string) (err error) {
	os.MkdirAll(filepath.Join(root, "site", "images"), 0755)
	os.MkdirAll(filepath.Join(root, "site", "javascript"), 0755)
	copyDirContents(root, "images")
	copyDirContents(root, "javascript")

	var pageFiles, postFiles []string
	pageFiles, err = relativePaths(root, "pages", "*.md")
	if err != err {
		fmt.Printf("Failed to find pages. Error: %s\n", err)
		return
	}
	postFiles, err = relativePaths(root, "posts", "*.md")
	if err != nil {
		fmt.Printf("Failed to find any .md files in posts subdir. Error: %s\n", err)
		return
	}
	var posts []*Post
	posts, err = postsMetadata(root, postFiles)
	if err != nil {
		fmt.Printf("Failed to get post metadata. Error: %s\n", err)
	}
	fmt.Printf("p: %v\n", posts)

	for _, pageFilepath := range pageFiles {
		templatedPage(root, pageFilepath, filepath.Base(pageFilepath))
	}

	for _, p := range posts {
		templatedPost(root, p)
	}

	createAtomFeed(filepath.Join(root, "site", "index.rss"), posts)
	return nil
}

func getCWD() (directory string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("OMG can't work out the CWD")
	}
	return dir
}

func main() {
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
				directory = getCWD()
			}
			checkDirectory(directory)
			generate(directory)
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

const aboutPage = `## About
This is the *about* page.
`

const indexPage = `## Index
This is the index page.
`

const examplePost = `---
layout: post
title: Example post
date: 2015-02-27 17:31
comments: true
categories: [tmux, cli, tips]
---

## Example blog post

Hi all this is my new blog isn't it the best no really.

It's got smart quotes and everything.
`

const layoutPjax = `
<article>{{.PageContent}}</article>
`

const layoutDefault = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
	<script src="http://ajax.googleapis.com/ajax/libs/jquery/1.5.2/jquery.min.js" type="text/javascript"></script>
	<script src="http://ajax.googleapis.com/ajax/libs/jqueryui/1.8.12/jquery-ui.min.js" type="text/javascript"></script>
	<script src="https://github.com/defunkt/jquery-pjax/raw/master/jquery.pjax.js" type="text/javascript"></script>
    <script src="/javascript/app.js" type="text/javascript"></script>
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
      <img src="http://d3nwyuy0nl342s.cloudfront.net/images/modules/facebox/loading.gif">
    </div>
    <div id="main">
      <article>
        {{.PageContent}}
      </article>
    </div>
  </body>
</html>
`

const exampleConfig = `---
site_name: My awesome blog
subtitle: Really great stuff
`
