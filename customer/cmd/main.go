package main

import (
	"net/http"
	"os"

	"reconcip.com.br/microservices/customer/pkg"
)

func main() {
	svc := pkg.NewService(pkg.NewValidator())
	httpHandler := pkg.MakeHTTPHandler(svc, os.Getenv("AUTH_SERVICE_URL"))
	http.ListenAndServe(":80", httpHandler)
}
