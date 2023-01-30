package models

// MailData holds email
type MailData struct {
	Name     string
	From     string
	Subject  string
	Content  string
	To       string
	Template string
}

// MailServer holds mail server parameters
type MailServer struct {
	Host     string
	Port     int
	Username string
	Password string
}
