package main

import (
	"io"
	"github.com/gomarkdown/markdown/ast"
	"strconv"
)

var headings = []Heading{}
var inHeading bool = false
var headingLevel int = 0

type Heading struct {
	Level int
	Title string
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if text, ok := node.(*ast.Link); ok {
		renderLinkHTMLExtension(w, text, entering)
		return ast.GoToNext, true
	}
	if heading, ok := node.(*ast.Heading); ok {
		renderHeading(w, heading, entering)
		return ast.GoToNext, true
	}
	if text, ok := node.(*ast.Text); ok {
		renderText(w, text)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

// Function to render links, renders normally except changes .md extension to .html extension
func renderLinkHTMLExtension(w io.Writer, p *ast.Link, entering bool) {
	// Format: <a href="https://www.nasa.gov/multimedia/imagegallery/iotd.html" target="_blank">NASA Image of the Day</a>
	if entering {
		var link string = string(p.Destination)
		if len(link) > 0 && link[len(link)-3:] == ".md" {
			link = link[:len(link)-3] + ".html"
		}
		io.WriteString(w, "<a href=\"" + link + "\" target=\"_self\">")
	} else {
		io.WriteString(w, "</a>")
	}
}

// Function to render heading including the heading ID for contents
func renderHeading(w io.Writer, p *ast.Heading, entering bool) {
	headingLevel = p.Level
	level := strconv.Itoa(p.Level)
	if entering {
		io.WriteString(w, "<h" + level + " id=" + p.HeadingID + ">")
		inHeading = true
	} else {
		io.WriteString(w, "</h" + level + ">")
		inHeading = false
	}
}

// Text rendering function, works normally except if the text is a heading, in which case appends the heading to the heading array
func renderText(w io.Writer, p *ast.Text) {
	var textString string = string(p.Literal)
	io.WriteString(w, textString)
	if inHeading {
		var curHeading Heading = Heading {
			Level: headingLevel,
			Title: textString,
		}
		headings = append(headings, curHeading)
	}
}
