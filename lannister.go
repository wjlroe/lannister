package main

import (
	 markdown "github.com/knieriem/markdown"
	"os"
	"io/ioutil"
	"io"
	"bufio"
	"fmt"
)

func markdown_parse(file_in io.Reader, file_out io.Writer) {
	b, _ := ioutil.ReadAll(file_in)

	doc := markdown.Parse(string(b), markdown.Extensions{Smart: true})

	w := bufio.NewWriter(file_out)
	doc.WriteHtml(w)
	w.Flush()
}

func createsite(site_dir String) {
	// Create the basic layout for a site
	abs_path = absolute(site_dir)
	os.Mkdir(abd_path, 0776)

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
		default:
			fmt.Printf("Command: %s not understood.\n", command)
		}
	} else {
		fmt.Println("no command specified!")
	}
}
