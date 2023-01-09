package main

import (
	"net/http"
	"os"

	"reconcip.com.br/microservices/supplier/pkg"
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

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	svc := pkg.NewService(validator, repository)
	svc = pkg.NewLoggingService(svc, logger)
}
