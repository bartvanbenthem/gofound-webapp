package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bartvanbenthem/gofound-webapp/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	blogposts      models.BlogPostModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	mailChan       chan models.MailData
	mailServer     models.MailServer
	mailAddress    string
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/gofound?parseTime=true", "MySQL data source name")
	smtphost := flag.String("smtp-host", "localhost", "Mail server hostname")
	smtpport := flag.Int("smtp-port", 1025, "Mail Server port")
	smtpuser := flag.String("smtp-user", "", "Mail Server username")
	smtppass := flag.String("smtp-password", "", "Mail Server password")
	mailaddr := flag.String("mail-address", "mail@gofound.nl", "Mail address")
	tlscert := flag.String("cert", "./tls/cert.pem", "Full path and TLS cert")
	tlskey := flag.String("key", "./tls/key.pem", "Full path and TLS key")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	sessionManager.Cookie.Secure = true

	// setting up mail channel
	mailChan := make(chan models.MailData)

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		blogposts:      &models.BlogPostModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		mailChan:       mailChan,
		mailServer: models.MailServer{
			Host:     *smtphost,
			Port:     *smtpport,
			Username: *smtpuser,
			Password: *smtppass,
		},
		mailAddress: *mailaddr,
	}

	// Mail Listener channel
	infoLog.Printf("starting mail listener\n")
	app.listenForMail()
	defer close(app.mailChan)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS(*tlscert, *tlskey)
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
