package main

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday/v2"
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"
)

func main() {
	input, err := os.ReadFile("/Users/santhoshreddy/Downloads/all_types_test.md")
	fmt.Println(err)
	output := blackfriday.Run([]byte(input), blackfriday.WithNoExtensions())
	err = os.WriteFile("all_types.html", output, 0644)
	fmt.Println(err)
	//doc, _ := html.Parse(strings.NewReader(string(output)))
	//convertHtmlToBlock(doc)
	//parse_html(doc)
	//nodes := FindAllInternal(doc)
	//jsonData, _ := ConvertHtmlToJson(nodes)
	//for _, node := range jsonData {
	//	fmt.Println("Toplevel =======> ", node.Text, node.Name)
	//	if len(node.Elements) > 0 && node.Name != "body" && node.Name != "" && node.Name != "head" {
	//		for _, Inode := range node.Elements {
	//			fmt.Println("inner level =====>", Inode.Text, Inode.Name)
	//		}
	//	}
	//}
}

//func parse_html(n *html.Node) {
//	if n.Type == html.TextNode {
//		if parent == "h1" {
//			fmt.Println(n.Data)
//		}
//	}
//	for c := n.FirstChild; c != nil; c = c.NextSibling {
//		if n.Type == html.ElementNode {
//			parent = n.Data
//		}
//		parse_html(c)
//	}
//}

func FindAllInternal(node *html.Node) []*html.Node {
	matched := []*html.Node{}
	matched = append(matched, node)
	for c := node.FirstChild; c != nil; c = c.NextSibling {

		if c.Type == html.ElementNode {
			//parent = c.Data
		}
		found := FindAllInternal(c)
		if len(found) > 0 {
			matched = append(matched, found...)
		}
	}
	return matched
}

// ConvertHtmlToJson the given HTML nodes into JSON content where each
// HTML node is represented by the JsonNode structure.
func ConvertHtmlToJson(nodes []*html.Node) ([]JsonNode, error) {
	rootJsonNodes := make([]JsonNode, len(nodes))

	for i, n := range nodes {
		rootJsonNodes[i].populateFrom(n)
	}
	return rootJsonNodes, nil
}

// JsonNode is a JSON-ready representation of an HTML node.
type JsonNode struct {
	// Name is the name/tag of the element
	Name string `json:"name,omitempty"`
	// Attributes contains the attributs of the element other than id, class, and href
	Attributes map[string]string `json:"attributes,omitempty"`
	// Class contains the class attribute of the element
	Class string `json:"class,omitempty"`
	// Id contains the id attribute of the element
	Id string `json:"id,omitempty"`
	// Href contains the href attribute of the element
	Href string `json:"href,omitempty"`
	// Text contains the inner text of the element
	Text string `json:"text,omitempty"`
	// Elements contains the child elements of the element
	Elements []JsonNode `json:"elements,omitempty"`
}

func (n *JsonNode) populateFrom(htmlNode *html.Node) *JsonNode {
	switch htmlNode.Type {
	case html.ElementNode:
		n.Name = htmlNode.Data
		break

	case html.TextNode:
		n.Name = htmlNode.Data
		break

	case html.DocumentNode:
		break

	default:
		log.Fatal("Given node needs to be an element or document")
	}

	var textBuffer bytes.Buffer

	if len(htmlNode.Attr) > 0 {
		n.Attributes = make(map[string]string)
		var a html.Attribute
		for _, a = range htmlNode.Attr {
			switch a.Key {
			case "class":
				n.Class = a.Val

			case "id":
				n.Id = a.Val

			case "href":
				n.Href = a.Val

			default:
				n.Attributes[a.Key] = a.Val
			}
		}
	}

	e := htmlNode.FirstChild
	for e != nil {
		switch e.Type {
		case html.TextNode:
			trimmed := strings.TrimSpace(e.Data)
			if len(trimmed) > 0 {
				// mimic HTML text normalizing
				if textBuffer.Len() > 0 {
					textBuffer.WriteString(" ")
				}
				textBuffer.WriteString(trimmed)
			}

		case html.ElementNode:
			if n.Elements == nil {
				n.Elements = make([]JsonNode, 0)
			}
			var jsonElemNode JsonNode
			jsonElemNode.populateFrom(e)
			n.Elements = append(n.Elements, jsonElemNode)
		}

		e = e.NextSibling
	}

	if textBuffer.Len() > 0 {
		n.Text = textBuffer.String()
	}

	return n
}

var parent string
var parentTracker = -1
var childTracker = 0
var sameP bool
var newP bool
var blockText string

func convertHtmlToBlock(node *html.Node) {
	if node.Data != "" && parent == "p" {
		if childTracker == parentTracker {
			blockText = node.Data
			fmt.Println(blockText)
		} else if childTracker != parentTracker && len(blockText) > 0 {
			blockText = ""
			sameP = false
			newP = false
			childTracker = parentTracker
			blockText = blockText + node.Data
			fmt.Println(blockText)
		}
	} else {
		if len(blockText) > 0 {
			blockText = ""
			sameP = false
			newP = false
			fmt.Println(blockText)
		}
		if node.Type == html.DocumentNode || node.Type == html.TextNode && len(node.Data) > 0 {
			blockText = ""
			sameP = false
			newP = false
			fmt.Println(blockText)
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		//fmt.Println(c.Data, c.Type)
		if parent == "ol" || parent == "ul" || parent == "h1" || parent == "h2" || parent == "h3" || parent == "h4" || parent == "h5" || parent == "h6" || parent == "p" || parent == "text" {
			convertHtmlToBlock(c)
		} else if c.Type == html.ElementNode && c.Data == "ol" {
			parent = "ol"
			convertHtmlToBlock(c)
			parent = ""
		} else if c.Type == html.ElementNode && c.Data == "li" {
			parent = "ul"
			convertHtmlToBlock(c)
			parent = ""
		} else if c.Type == html.ElementNode && (c.Data == "h1" || c.Data == "h2" || c.Data == "h3" || c.Data == "h4" || c.Data == "h5" || c.Data == "h6") {
			parent = c.Data
			convertHtmlToBlock(c)
			parent = ""
		} else if c.Type == html.ElementNode && c.Data == "p" {
			parentTracker++
			parent = c.Data
			sameP = true
			newP = true
			convertHtmlToBlock(c)
			sameP = true
			parent = ""
		} else if c.Type == html.ElementNode && c.Data == "li" {
			convertHtmlToBlock(c)
		} else {
			parent = ""
			convertHtmlToBlock(c)
		}
	}
}
