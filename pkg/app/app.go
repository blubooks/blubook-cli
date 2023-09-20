package app

import (
	"log"
	"os"
	"text/template"
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

type App struct {
}

func New() *App {
	return &App{}
}

func Build() error {
	menu := BookNavi()

	tmpl, err := template.ParseGlob("data/layout/default/*")
	if err != nil {
		return err
	}

	for _, s := range menu.Pages {

		log.Printf("%+v", s)

		if s.Link != nil {

			if err != nil {
				return err
			}
			log.Println("public/" + *s.Link + ".html")
			file, err := os.Create("public/" + *s.Link + ".html")
			if err != nil {
				return err
			}

			err = tmpl.ExecuteTemplate(file, "index.html", menu)
			if err != nil {
				return err
			}

		}

	}

	return nil

}
