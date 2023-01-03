package main

import (
	"net/http"

	"api.example.com/customer/pkg"
)

func main() {
	svc := pkg.NewService(pkg.NewValidator())
	http.ListenAndServe(":80", pkg.MakeHTTPHandler(svc))
}
