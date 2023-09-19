package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Menu struct {
	Title string
}
type Page struct {
	Title string
	Link  string
	Page  []Page
}
type MenuEntry struct {
	Type       int
	GroupTitle string
	Pages      []Page
}

func main() {
	//fmt.Println("hello world")

	_, err := os.ReadFile("test.md")
	check(err)
	//fmt.Print(string(dat))

	f, err := os.Open("test.md")
	check(err)

	b1 := make([]byte, 1000)
	_, err = f.Read(b1)
	check(err)
	//fmt.Printf("%d bytes: %s\n", n1, string(b1[:n1]))

	source := b1

	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		panic(err)
	}

	/*

		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithXHTML(),
			),
		)

		doc := md.Parser().Parse(text.NewReader(b1))

		ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {

			//if entering {
			fmt.Printf("WAlk: %+v\n", n.Kind())

			//}

			return ast.WalkContinue, nil
		})
	*/

	doc := goldmark.DefaultParser().Parse(text.NewReader([]byte(source)))
	listLevel := 0

	var menu Menu
	menu.Title = "test"
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error

		if entering {

			fmt.Println("   ")
			fmt.Println(n.Kind(), entering, string(n.Text([]byte(source))))

			if n.Kind() == ast.KindHeading {
				firstH := n.(*ast.Heading)
				fmt.Printf("Level  %+v\n", firstH.Level)
			}
			if n.Kind() == ast.KindListItem {
				firstH := n.(*ast.ListItem)
				firstH.SetAttribute([]byte("level"), 1)
				//fmt.Printf("- parent  %+v\n", firstH.Parent())
			}
			if n.Kind() == ast.KindList {
				listLevel = listLevel + 1

				firstH := n.(*ast.List)
				firstH.SetAttribute([]byte("level"), listLevel)
				fmt.Printf("ListBlock  %+v\n", firstH.Parent().Kind())
				test, _ := firstH.Attribute([]byte("level"))
				fmt.Printf("Listlevel  %+v\n", test)
			}

		} else {
			if n.Kind() == ast.KindList {
				listLevel = listLevel - 1
			}
		}

		return s, err
	})

}
