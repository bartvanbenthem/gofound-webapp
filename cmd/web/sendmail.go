package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
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
		log.Printf("Error: %s\n", err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./templates/email/%s", m.Template))
		if err != nil {
			log.Printf("Error: %s\n", err)
		}
		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		log.Printf("Error: %s\n", err)
	} else {
		log.Printf("Email sent From: %s To: %s Subject: %s\n",
			m.From, m.To, m.Subject)
	}
}
