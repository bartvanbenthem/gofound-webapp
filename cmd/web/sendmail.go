package main

import (
	"log"
	"time"

	"github.com/bartvanbenthem/gofound-webapp/internal/models"
	mail "github.com/xhit/go-simple-mail"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMessage(msg)
		}
	}()
}

func sendMessage(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.Username = ""
	server.Password = ""

	client, err := server.Connect()
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		log.Printf("Error: %s\n", err)
	} else {
		log.Printf("Email sent From: %s To: %s Subject: %s\n",
			m.From, m.To, m.Subject)
	}
}
