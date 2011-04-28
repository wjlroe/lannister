package main

import (
	 markdown "github.com/knieriem/markdown"
	"os"
	"io/ioutil"
	"io"
//	"bufio"
	"fmt"
	"path/filepath"
	"bytes"
	"template"
	"strings"
)

type Page struct {
	PageContent string
}

func is_dir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDirectory()
}

func check_cwd() {
	if !is_dir("pages") || !is_dir("layouts") || !is_dir("site") {
		fmt.Println("Current directory is not a Lannister site dir.")
		fmt.Println("To be used, there should be ./pages/, ./layouts/, ./site/")
		fmt.Println("Please use ./lannister createsite ./dir/to/create/as/site")
		os.Exit(1)
	}
}

// TODO: Return the markdown doc so it can be used multiple times instead of parsing it again
func markdown_parse(file_in io.Reader) *markdown.Doc {
	b, _ := ioutil.ReadAll(file_in)

	doc := markdown.Parse(string(b), markdown.Extensions{Smart: true})

	//w := bufio.NewWriter(file_out)
	//doc.WriteHtml(w)
	//w.Flush()
	return doc
}

func createsite(site_dir string) {
	// Create the basic layout for a site
	//abs_path := absolute(site_dir)
	//os.Mkdir(abd_path, 0776)

}

func generate() os.Error {
	// list the files to be templated
	// process them - markdown
	// run each through each layout file (TODO: make configurable)
	// save each output file in the site directory
	var page_files []string
	var err os.Error
	page_files, err = filepath.Glob("./pages/*.md")
	if err != nil {
		fmt.Printf("Failed to find any .md files in pages subdir. Error: %s\n", err.String())
		return err
	}

	var template_files []string
	template_files, err = filepath.Glob("./layouts/*.html")
	if err != nil {
		fmt.Printf("Failed to find any applicable layout files. Error: %s\n", err.String())
		return err
	}
	var fmap = template.FormatterMap{
		"" : template.StringFormatter,
		"html": template.HTMLFormatter,
	}
	//templates := make([]template.Template, len(template_files))
	templates := map[string] *template.Template{}
	for _,t := range template_files {
		templates[filepath.Base(t)] = template.MustParseFile(t, fmap)
	}

	for _,page_filepath := range page_files {
		in_fd, err := os.Open(page_filepath)
		if err != nil {
			fmt.Printf("Failed to open file: %s, error: %s\n", page_filepath, err.String())
			return err
		}
		defer in_fd.Close()
		doc := markdown_parse(in_fd)
		buffer := bytes.NewBufferString("")
		doc.WriteHtml(buffer)

		page := &Page{PageContent: buffer.String()}

		page_filename := filepath.Base(page_filepath)
		page_ext := filepath.Ext(page_filename)
		fmt.Printf("Page ext: %s\n", page_ext)
		page_filename = strings.Replace(page_filename, page_ext, "", -1)
		for tpl_filename, template := range templates {
			filename := strings.Replace(tpl_filename, "filename", page_filename, -1)
			fmt.Printf("tpl_filename: %s, page_filename: %s, Output page filename: %s\n", tpl_filename, page_filename, filename)
			filepath := filepath.Join("site", filename)
			fmt.Printf("Going to save templated page: %s as file: %s\n", page_filename, filepath)
			out_fd, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				fmt.Printf("Failed to open file: %s for writing, error: %s\n", filepath, err.String())
				return err
			}
			defer out_fd.Close()
			template.Execute(out_fd, page)
			out_fd.Sync()
		}
	}
	return nil
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
			check_cwd()
			generate()
		default:
			fmt.Printf("Command: %s not understood.\n", command)
		}
	} else {
		fmt.Println("no command specified!")
	}
}
