package main

import (
	"net/http"

	"reconcip.com.br/microservices/auth/pkg"
)

func main() {
	svc := pkg.NewService(pkg.NewTokenGenerator())
	http.ListenAndServe(":80", pkg.NewHTTPHandler(svc))
}
