package main

import (
	"time"

	"github.com/bartvanbenthem/gofound-webapp/internal/models"
	mail "github.com/xhit/go-simple-mail"
)

func (app *application) listenForMail() {
	go func() {
		for {
			msg := <-app.mailChan
			app.sendMessage(msg)
		}
	}()
}

func (app *application) sendMessage(m models.MailData) {

	server := mail.NewSMTPClient()
	server.Host = app.mailServer.Host
	server.Port = app.mailServer.Port
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.Username = app.mailServer.Username
	server.Password = app.mailServer.Password

	client, err := server.Connect()
	if err != nil {
		app.errorLog.Printf("%s", err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		app.errorLog.Printf("%s", err)
	} else {
		app.infoLog.Printf("Email sent From: %s To: %s\n", m.From, m.To)
	}
}
