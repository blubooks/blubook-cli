package blubookcli

// test
import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

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
	Type  int    `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
	Page  *Page  `json:"page,omitempty"`
	Pages []Page `json:"pages,omitempty"`
	Set   bool   `json:"-"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func list(node ast.Node, initLevel int, page *Page, source *[]byte) {
	level := initLevel
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)

		if entering {
			if n.Kind() == ast.KindList {
				level = level + 1

			}
			if n.Kind() == ast.KindListItem {
				if level == initLevel+1 {
					pg := Page{}
					pg.Title, pg.Link = listitemlink(n.FirstChild(), source)

					list(n, level, &pg, source)

					page.Pages = append(page.Pages, pg)

				}
			}
		} else {
			if n.Kind() == ast.KindList {
				level = level - 1
			}
		}
		var err error
		return s, err
	})
}

func listitemlink(node ast.Node, source *[]byte) (text string, link string) {
	l_text := ""
	l_link := ""
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)

		if entering {
			if n.Kind() == ast.KindLink {
				l := n.(*ast.Link)
				l_text = string(n.Text([]byte(*source)))
				l_link = string(l.Destination)

			}
		}
		var err error
		return s, err
	})

	if l_text == "" {
		l_text = string(node.FirstChild().Text([]byte(*source)))
	}

	return l_text, l_link

}

func BookNavi() *Menu {

	source, err := os.ReadFile("content/SUMMARY.md")
	check(err)

	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		panic(err)
	}

	doc := goldmark.DefaultParser().Parse(text.NewReader([]byte(source)))
	listLevel := 0

	var menu Menu
	var page Page
	var entry MenuEntry
	menu.Title = "TITLE"
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
						Set:   true,
						Type:  1,
						Title: string(n.Text([]byte(source))),
					}
				}

			} else if n.Kind() == ast.KindThematicBreak {
				if entry.Set {
					menu.Entries = append(menu.Entries, entry)
					entry = MenuEntry{}

				}
			} else if n.Kind() == ast.KindList {
				listLevel = listLevel + 1

			} else if n.Kind() == ast.KindListItem {

				if listLevel == 1 {

					if entry.Type == 1 {
						page = Page{}
						page.Title, page.Link = listitemlink(n.FirstChild(), &source)
						list(n, 1, &page, &source)

						entry.Pages = append(entry.Pages, page)
					} else {
						if entry.Set {
							menu.Entries = append(menu.Entries, entry)
						}

						pg := Page{}
						pg.Title, pg.Link = listitemlink(n, &source)
						entry = MenuEntry{
							Set:  true,
							Type: 3,
							Page: &pg,
						}
						menu.Entries = append(menu.Entries, entry)
						entry = MenuEntry{}

					}
				}
			}

		} else {
			if n.Kind() == ast.KindList {
				listLevel = listLevel - 1
			} else if n.Kind() == ast.KindDocument {
				if entry.Set {
					menu.Entries = append(menu.Entries, entry)
				}
			}
		}

		return s, err
	})

	return &menu
}
