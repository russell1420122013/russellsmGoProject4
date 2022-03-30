package main

import (
	"log"
	"net/http"
)

const (
	AddSrv       = ":8000"
	TemplatesDir = "."
)

func main() {
	log.Printf("listening on %q...", AddSrv)
	fileSrv := http.FileServer(http.Dir(TemplatesDir))
	if err := http.ListenAndServe(AddSrv, fileSrv); err != nil {
		log.Fatal(err)
	}
}
