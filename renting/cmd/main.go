package main

import (
	"net/http"
	"os"

	"reconcip.com.br/microservices/renting/pkg"
)

func main() {
	repository, err := pkg.NewMongoRepository(
		os.Getenv("MONGODB_URL"),
		os.Getenv("MONGODB_USER"),
		os.Getenv("MONGODB_PASSWORD"),
		os.Getenv("MONGODB_DATABASE"),
	)

	if err != nil {
		panic(err)
	}

	validator := pkg.NewValidator([]pkg.ValidationRule{
	})

	svc := pkg.NewService(validator, repository)
	endpoints := pkg.CreateEndpoints(svc)
	http.ListenAndServe(":80", pkg.NewHTTPServer(endpoints))
}
