package app

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wellington/go-libsass"
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
	PublicURL  string
}

func copyDir(source, destination string) error {
	var err error = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, source, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(destination, relPath), 0755)
		} else {
			var data, err1 = ioutil.ReadFile(filepath.Join(source, relPath))
			if err1 != nil {
				return err1
			}
			return ioutil.WriteFile(filepath.Join(destination, relPath), data, 0777)
		}
	})
	return err
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

func Build(publicURL, layoutPath string, publicPath string, dataPath string) (error, []string) {
	if layoutPath == "" {
		layoutPath = "data/layout/default/"
	}
	if publicPath == "" {
		publicPath = "public/"
	}
	if dataPath == "" {
		dataPath = "data/content/"
	}
	if publicURL == "" {
		publicURL = "http://localhost/public/"
	}

	dataPath = setLastLash(dataPath)
	layoutPath = setLastLash(layoutPath)
	publicPath = setLastLash(publicPath)
	publicURL = setLastLash(publicURL)

	menu := BookNavi()
	var out []string
	var data TemplateData
	data.Menu = menu
	data.PublicPath = publicPath
	data.ActivePath = layoutPath
	data.DataPath = dataPath
	data.PublicURL = publicURL

	tmpl, err := template.New("").Funcs(funcMap).ParseGlob(layoutPath + "*.html")
	if err != nil {
		return err, []string{err.Error()}
	}
	RemoveGlob(publicPath + "*")

	copyDir(layoutPath+"static", publicPath)
	_ = os.MkdirAll(publicPath+".data/assets", os.ModePerm)
	copyDir(dataPath+".data/assets", publicPath+".data/assets")

	file, err := os.Create(publicPath + "index.html")
	if err != nil {
		return err, []string{err.Error()}
	}

	err, data.Content = loadMarkdown(dataPath + "README.md")
	if err != nil {
		return err, []string{err.Error()}
	}
	data.ActivePath = publicPath

	err = tmpl.ExecuteTemplate(file, "index.html", data)
	if err != nil {
		log.Printf("Error in Build() -> tmpl.ExecuteTemplate(): %v", err)
		out = append(out, err.Error())
	}

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

	styleFile := layoutPath + "assets/style/style.scss"
	if _, err := os.Stat(styleFile); err == nil {
		r, err := os.Open(styleFile)
		if err != nil {
			log.Println(err)
		}
		var styleBuffer bytes.Buffer

		comp, err := libsass.New(&styleBuffer, r)
		if err != nil {
			out = append(out, "Style: "+err.Error())
		}
		// configure @import paths
		includePaths := []string{layoutPath + "assets/style/partials"}
		err = comp.Option(libsass.IncludePaths(includePaths))
		if err != nil {
			out = append(out, "Style Options: "+err.Error())
		}

		if err := comp.Run(); err != nil {
			out = append(out, "Style Compiler: "+err.Error())

		}

		err = os.WriteFile(publicPath+"style.css", styleBuffer.Bytes(), os.ModePerm)
	}
	return nil, out

}
