package main

import (
	"net/http"

	"api.example.com/auth/pkg"
)

func main() {
	svc := pkg.NewService()
	http.ListenAndServe(":80", pkg.NewHTTPHandler(svc))
}
