package main

import (
	"bufio"
	"bytes"
	"encoding/json"
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
	Title   string      `json:"title,omitempty"`
	Entries []MenuEntry `json:"entries,omitempty"`
}
type Page struct {
	Title string `json:"title,omitempty"`
	Link  string `json:"link,omitempty"`
	Pages []Page `json:"pages,omitempty"`
}
type MenuEntry struct {
	Type       int    `json:"type,omitempty"`
	GroupTitle string `json:"groupTitle,omitempty"`
	Page       *Page  `json:"page,omitempty"`
	Pages      []Page `json:"pages,omitempty"`
	Set        bool   `json:"-"`
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
	var entry MenuEntry
	menu.Title = "test"
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error

		if entering {

			if n.Kind() == ast.KindHeading {

				h := n.(*ast.Heading)
				if h.Level == 2 {

					if entry.Set {
						menu.Entries = append(menu.Entries, entry)
					}

					entry = MenuEntry{
						Set:        true,
						Type:       1,
						GroupTitle: string(n.Text([]byte(source))),
					}
				}

			} else if n.Kind() == ast.KindList {
				listLevel = listLevel + 1
			} else if n.Kind() == ast.KindListItem {

				if listLevel == 1 {
					page := Page{
						Title: string(n.Text([]byte(source))),
					}

					if entry.Type != 1 {
						entry = MenuEntry{
							Set:  true,
							Type: 3,
						}
						entry.Page = &page
						menu.Entries = append(menu.Entries, entry)
						entry = MenuEntry{}
					} else {
						entry.Pages = append(entry.Pages, page)

					}

				}
			}

		} else {
			if n.Kind() == ast.KindList {
				listLevel = listLevel - 1
			} else if n.Kind() == ast.KindThematicBreak {
				menu.Entries = append(menu.Entries, entry)
				entry = MenuEntry{}
			}
		}
		/*
			if entering {

				fmt.Println("   ")
				fmt.Println(n.Kind(), entering, string(n.Text([]byte(source))))

				if n.Kind() == ast.KindHeading {

					h := n.(*ast.Heading)
					if h.Level == 2 {

						if entry.Type == 1 {
							menu.Entries = append(menu.Entries, entry)
							entry = MenuEntry{
								Type: 3,
							}
						}

						entry = MenuEntry{
							Type:       1,
							GroupTitle: string(n.Text([]byte(source))),
						}
					}

				} else if n.Kind() == ast.KindListItem {
					firstH := n.(*ast.ListItem)
					firstH.SetAttribute([]byte("level"), 1)
					page := Page{
						Title: string(n.Text([]byte(source))),
					}
					entry.Pages = append(entry.Pages, page)
				} else if n.Kind() == ast.KindList {
					listLevel = listLevel + 1
					if listLevel == 1 {
						if entry.Type == 0 {
							entry = MenuEntry{
								Type: 2,
							}
						}
					}
				} else if n.Kind() == ast.KindThematicBreak {
					listLevel = listLevel + 1
					menu.Entries = append(menu.Entries, entry)
					entry = MenuEntry{
						Type: 3,
					}

				}

			} else {
				if n.Kind() == ast.KindDocument {
					menu.Entries = append(menu.Entries, entry)

				}
			}
		*/

		return s, err
	})

	//fmt.Printf("- parent  %+v\n", menu)
	b, err := json.Marshal(menu)
	var out bytes.Buffer

	json.Indent(&out, b, "", "  ")

	if err != nil {
		fmt.Println(err)
		return
	}

	f, err = os.Create("test.json")
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)
	n4, err := w.WriteString(string(out.Bytes()))
	check(err)
	fmt.Printf("wrote %d bytes\n", n4)
	w.Flush()

}
