package blubookcli

// test
import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const (
	TypeGroup    int = 1
	TypeMenuItem int = 2
	TypeLink     int = 3
)

type Menu struct {
	Title string `json:"title,omitempty"`
	Pages []Page `json:"pages,omitempty"`
}
type Page struct {
	Set        bool   `json:"-"`
	Parent     *Page  `json:"-"`
	ParentLink string `json:"parent,omitempty"`
	Level      int    `json:"level,omitempty"`
	Type       int    `json:"type,omitempty"`
	Title      string `json:"title,omitempty"`
	Link       string `json:"link,omitempty"`
	Pages      []Page `json:"pages,omitempty"`
}

/*
type MenuEntry struct {
	Type  int    `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
	Pages []Page `json:"pages,omitempty"`
	Set   bool   `json:"-"`
}
*/

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
					pg.Type = TypeLink
					pg.Level = level
					pg.Parent = page
					pg.ParentLink = page.Link
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

	source, err := os.ReadFile("data/content/SUMMARY.md")
	check(err)

	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		panic(err)
	}

	doc := goldmark.DefaultParser().Parse(text.NewReader([]byte(source)))
	listLevel := 0

	var menu Menu

	var entry Page
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error

		if entering {

			if n.Kind() == ast.KindHeading {

				h := n.(*ast.Heading)
				if h.Level == 1 && menu.Title == "" {
					menu.Title = string(n.Text([]byte(source)))
				} else if h.Level == 2 {

					if entry.Set {
						menu.Pages = append(menu.Pages, entry)
					}

					entry = Page{
						Set:   true,
						Level: listLevel,
						Type:  TypeGroup,
						Title: string(n.Text([]byte(source))),
					}
				}

			} else if n.Kind() == ast.KindThematicBreak {
				if entry.Set {
					menu.Pages = append(menu.Pages, entry)
					entry = Page{}
				}
			} else if n.Kind() == ast.KindList {
				listLevel = listLevel + 1

			} else if n.Kind() == ast.KindListItem {

				if listLevel == 1 {

					pg := Page{}
					pg.Set = true
					pg.Type = TypeMenuItem
					pg.Level = listLevel
					pg.Title, pg.Link = listitemlink(n.FirstChild(), &source)

					if entry.Type == 1 {
						list(n, 1, &pg, &source)
						entry.Pages = append(entry.Pages, pg)
					} else {
						menu.Pages = append(menu.Pages, pg)
					}
				}
			}

		} else {
			if n.Kind() == ast.KindList {
				listLevel = listLevel - 1
			} else if n.Kind() == ast.KindDocument {
				if entry.Set {
					menu.Pages = append(menu.Pages, entry)
				}
			}
		}

		return s, err
	})

	return &menu
}
