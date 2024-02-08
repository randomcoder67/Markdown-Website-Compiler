package main

import (
	"os"
	"io"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	TYPE_FILE int = 0
	TYPE_DIR int = 1
	TYPE_NOTHING int = 2
)

// Check if a given path points to a file, a directory, or nothing
func isFileOrDir(filename string) int {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return TYPE_NOTHING
	} else if fileInfo.IsDir() {
		return TYPE_DIR
	}
	return TYPE_FILE
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	//htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions {
		Flags: html.CommonFlags | html.HrefTargetBlank,
		RenderNodeHook: myRenderHook,
	}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func renderFile(inputFile string, outputFile string, level int) {
	// Open file
	mds, _ := os.ReadFile(inputFile)
	// Convert markdown to html
	md := []byte(mds)
	html := mdToHTML(md)
	
	// Write html to file
	var cssInclude string = "<link rel=\"stylesheet\" type=\"text/css\" href=\"" + strings.Repeat("../", level) + "style.css\"/>"
	err := os.WriteFile(outputFile, []byte(cssInclude + string(html)), 0666)
	if err != nil {
		panic(err)
	}
}

func renderFolder(dirName string, outputDir string) {
	// Check dirName is a directory
	if isFileOrDir(dirName) != TYPE_DIR {
		fmt.Println("Error, not a directory")
		os.Exit(1)
	}
	
	// Check outputDir has a slash
	if outputDir[len(outputDir)-1] != '/' {
		outputDir = outputDir + "/"
	}
	// Make the "new root" directory
	if isFileOrDir(outputDir) != TYPE_DIR {
		os.Mkdir(outputDir, os.ModePerm)
	}
	// Copy style.css there
	homeDir, _ := os.UserHomeDir()
	source, err := os.Open(homeDir + CSS_FILE_LOCATION)
	if err != nil {
		panic(err)
	}
	// Open destination
	destination, err := os.Create(outputDir + "style.css")
	if err != nil {
		panic(err)
	}
	// Copy source to destination
	_, err = io.Copy(destination, source)
	if err != nil {
		panic(err)
	}
	
	// Recursive function to go through given folder and subfolders
	var listInDir func (string, int)
	listInDir = func(dirNameCur string, level int) {
		fmt.Println(level)
		files, err := ioutil.ReadDir(dirNameCur)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			if file.IsDir() { // If file is a directory, make it in new location and go a level deeper
				if isFileOrDir(strings.ReplaceAll(outputDir + dirNameCur + "/" + file.Name(), dirName, "")) != TYPE_DIR {
					err := os.Mkdir(strings.ReplaceAll(outputDir + dirNameCur + "/" + file.Name(), dirName, ""), os.ModePerm)
					if err != nil {
						panic(err)
					}
				}
				listInDir(dirNameCur + "/" + file.Name(), level+1)
			} else { // If file is a file, add it to allFiles
				// Get new file name
				var newFileName string = strings.ReplaceAll(outputDir + dirNameCur + "/" + file.Name(), dirName, "")
				// If file is a markdown file
				if newFileName[len(newFileName)-3:] == ".md" {
					// Change extension to html
					newFileName = newFileName[:len(newFileName)-3] + ".html"
					//fmt.Println(dirNameCur + "/" + file.Name(), newFileName)
					// Render file
					renderFile(dirNameCur + "/" + file.Name(), newFileName, level)
				} else { // If file is not a markdown file, just copy
					// Open source
					source, err := os.Open(dirNameCur + "/" + file.Name())
					if err != nil {
						panic(err)
					}
					// Open destination
					destination, err := os.Create(newFileName)
					if err != nil {
						panic(err)
					}
					// Copy source to destination
					_, err = io.Copy(destination, source)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
	listInDir(dirName[:len(dirName)-1], 0)
}
