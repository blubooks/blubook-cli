package app

import (
	"fmt"
	"log"
	"net/http"
)

func (app *App) HanlderHealth(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Write the status code using w.WriteHeader
	res.WriteHeader(http.StatusOK)

	// Write the body text using w.Write
	res.Write([]byte("OK"))
}

func (app *App) HandleIndex(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	err, warn := Build("http://localhost:3020/public/", "", "", "")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		log.Printf("Fehler: %v", err)
		fmt.Fprintf(res, `{"error.message": "%v"}`, appErr)
		return
	}

	if warn != nil {
		fmt.Fprintf(res, `{"warn": %v}`, warn)

	}

}
