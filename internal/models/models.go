package models

type TestForm struct {
	Name  string
	Email string
}

// MailData holds email
type MailData struct {
	Name    string
	From    string
	Subject string
	Content string
	To      string
}
