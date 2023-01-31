# gofound-webapp
Project template for creating web applications in Go, with example database & content.
* Routing & middleware
* HTTP session management
* User registration & authentication
* Database backend (SQL)
* Templating & forms

### Used packages
* github.com/alexedwards/scs/mysqlstore
* github.com/alexedwards/scs/v2
* github.com/go-playground/form/v4
* github.com/go-sql-driver/mysql
* github.com/julienschmidt/httprouter
* github.com/justinas/alice
* github.com/justinas/nosurf
* golang.org/x/crypto


## Clone repo
```bash 
git clone https://github.com/bartvanbenthem/gofound-webapp.git
cd gofound-webapp
```

## Start MySQL Database
```shell
cd project
docker-compose up -d
cd ..
```

## Build & start the WebServer
```bash
go test -vet=off -v ./cmd/web/

go build -o ./bin/gofoundweb ./cmd/web/

./bin/gofoundweb --addr=":4000" \
                 --dsn="web:pass@/gofound?parseTime=true" \
                 --smtp-host="localhost" \
                 --smtp-port="1025" \
                 --smtp-user="" \
                 --smtp-password="" \
                 --mail-address="mail@gofound.nl" \
                 --cert="./tls/cert.pem" \
                 --key="./tls/key.pem"
```

## Test SendMail
Start MailHog to test mail capabillities
```shell
$ go get github.com/mailhog/MailHog
$ ~/go/bin/MailHog
```

## Test login
Use the provisioned user:
* Username: admin@gofound.nl 
* Password: administrator


## create self signed certificates
```bash
cd tls
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```