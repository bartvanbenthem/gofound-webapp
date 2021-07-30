package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bartvanbenthem/gofound-web/pkg/config"
	"github.com/bartvanbenthem/gofound-web/pkg/handlers"
	"github.com/bartvanbenthem/gofound-web/pkg/render"
)

const portNumber = ":8080"

// main is the main function
func main() {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
