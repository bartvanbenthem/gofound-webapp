package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bartvanbenthem/gofound-webapp/internal/config"
	"github.com/bartvanbenthem/gofound-webapp/internal/handlers"
	"github.com/bartvanbenthem/gofound-webapp/internal/models"
	"github.com/bartvanbenthem/gofound-webapp/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main function
func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	// Mail Listener channel
	fmt.Printf("starting mail listener\n")
	listenForMail()
	defer close(app.MailChan)

	// starting application
	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// declare what the session can consist of
	gob.Register(models.TestForm{})
	gob.Register(models.MailData{})

	// setting up mail channel
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = false

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	return nil
}
