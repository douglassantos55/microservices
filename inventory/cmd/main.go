package main

import (
	"net/http"
	"os"

	"reconcip.com.br/microservices/inventory/pkg"
)

func main() {
	validator := pkg.NewValidator()

	repository, err := pkg.NewMongoRepository(
		os.Getenv("MONGODB_URL"),
		os.Getenv("MONGODB_USER"),
		os.Getenv("MONGODB_PASSWORD"),
		os.Getenv("MONGODB_DATABASE"),
	)

	if err != nil {
		panic(err)
	}

	endpoints := pkg.NewSet(pkg.NewService(validator, repository))
	http.ListenAndServe(":80", pkg.NewHTTPHandler(endpoints))
}
