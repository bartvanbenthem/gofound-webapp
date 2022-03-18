# gofound-webapp
A project template for creating an web application in Go, with example content.


### Used packages
* github.com/alexedwards/scs/v2 v2.4.0 
* github.com/asaskevich/govalidator v11.0.0
* github.com/go-chi/chi v1.5.4
* github.com/justinas/nosurf v1.1.1
* github.com/xhit/go-simple-mail v2.2.2

## Prerequisites
Start MailHog to test mail capabillities
```shell
$ go get github.com/mailhog/MailHog
$ ~/go/bin/MailHog
```

## Start the WebServer
```shell
$ ./build/bin/webserver
```
```bash
starting mail listener
Staring application on port :8080
```
