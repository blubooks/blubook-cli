package app

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
)

var version = "0.0.0"

const (
	appErr                    = "app error"
	appErrCreationFailure     = "error createn failure"
	appErrDataAccessFailure   = "data access failure"
	appErrJsonCreationFailure = "json creation failure"

	appErrDataCreationFailure = "data creation failure"
	appErrFormDecodingFailure = "form decoding failure"

	appErrDataUpdateFailure      = "data update failure"
	appErrFormErrResponseFailure = "form error response failure"
)

var funcMap = template.FuncMap{
	"partial": partial,
}

type App struct {
}

func New() *App {
	return &App{}
}

type TemplateData struct {
	Menu       *Menu
	Page       *Page
	ActivePath string
	Content    string
	PublicPath string
	DataPath   string
	LayoutPath string
}

func partial(name string, data any) string {

	filename := filepath.Base(name)

	tmpl, err := template.New(filename).ParseFiles("data/layout/default/partial/" + name)

	if err != nil {
		return "Partial \"" + name + "\" not Found!"
	}

	var b bytes.Buffer

	tmpl.Execute(&b, data)

	result := b.String()

	return result
}

func loadMarkdown(filename string) (error, string) {
	source, err := os.ReadFile(filename)
	if err != nil {
		return err, ""
	}
	var buf bytes.Buffer

	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	if err := markdown.Convert(source, &buf); err != nil {
		return err, ""
	}
	return nil, string(buf.Bytes())
}

func RemoveGlob(path string) (err error) {
	contents, err := filepath.Glob(path)
	if err != nil {
		return
	}
	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return
		}
	}
	return
}

func createHtmlFiles(tmpl *template.Template, page Page, out *[]string, data TemplateData) {

	if page.DataLink == nil || page.Link == nil {
		return
	}

	publicFilePath := data.PublicPath + *page.Link
	_ = os.MkdirAll(publicFilePath, os.ModePerm)

	var file *os.File
	var err error

	publicFilePath = setLastLash(publicFilePath)

	file, err = os.Create(publicFilePath + "index.html")
	if err != nil {
		log.Printf("Error in createHtmlFiles -> os.Create(): %v ", err)
		*out = append(*out, err.Error())
	}

	err = tmpl.ExecuteTemplate(file, "index.html", data)
	if err != nil {
		log.Printf("Error in createHtmlFiles -> tmpl.ExecuteTemplate(): %v", err)
		*out = append(*out, err.Error())
	}

}

func page(pg []Page, tmpl *template.Template, out *[]string, data TemplateData) {

	for _, s := range pg {

		if s.Link != nil && s.DataLink != nil {

			var err error
			err, data.Content = loadMarkdown(data.DataPath + *s.DataLink)
			data.ActivePath = *s.Link
			data.Page = &s

			if err != nil {
				log.Printf("Error in err, page() -> loadMarkdown(): %v", err)
				*out = append(*out, err.Error())
			}

			createHtmlFiles(tmpl, s, out, data)

		}
		if len(s.Pages) > 0 {
			page(s.Pages, tmpl, out, data)
		}
	}
}

func setLastLash(text string) string {
	last := text[len(text)-1:]

	if last != "/" {
		return text + "/"

	}
	return text

}

func Build(layoutPath string, publicPath string, dataPath string) (error, []string) {
	if layoutPath == "" {
		layoutPath = "data/layout/default/"
	}
	if publicPath == "" {
		publicPath = "public/"
	}
	if dataPath == "" {
		dataPath = "data/content/"
	}

	dataPath = setLastLash(dataPath)
	layoutPath = setLastLash(layoutPath)
	publicPath = setLastLash(publicPath)

	menu := BookNavi()
	var out []string
	var data TemplateData
	data.Menu = menu
	data.ActivePath = publicPath
	data.PublicPath = publicPath
	data.ActivePath = layoutPath
	data.DataPath = dataPath

	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(layoutPath + "*.html")
	if err != nil {
		return err, []string{err.Error()}
	}

	RemoveGlob(publicPath + "*")

	file, err := os.Create(publicPath + "index.html")
	if err != nil {
		return err, []string{err.Error()}
	}

	err, data.Content = loadMarkdown(dataPath + "README.md")
	if err != nil {
		return err, []string{err.Error()}
	}

	err = tmpl.ExecuteTemplate(file, "index.html", data)
	if err != nil {
		log.Printf("Error in Build() -> tmpl.ExecuteTemplate(): %v", err)
		out = append(out, err.Error())
	}
	data.Content = ""

	for _, s := range menu.Pages {

		if s.DataLink != nil {

			err, data.Content = loadMarkdown(dataPath + *s.DataLink)
			data.ActivePath = *s.Link
			data.Page = &s
			if err != nil {
				log.Printf("Error in Build() -> loadMarkdown(): %v", err)
				out = append(out, err.Error())
			}

			createHtmlFiles(tmpl, s, &out, data)
		}

		if len(s.Pages) > 0 {
			page(s.Pages, tmpl, &out, data)
		}

	}

	return nil, out

}
