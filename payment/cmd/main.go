package main

import (
	"net/http"
	"os"

	"github.com/go-kit/log"
	"reconcip.com.br/microservices/payment/pkg"
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

	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	validator := pkg.NewValidator([]pkg.ValidationRule{
		pkg.NewPaymentTypeRule(repository),
	})
	svc := pkg.NewService(validator, repository)
	svc = pkg.NewLoggingService(svc, logger)

	endpoints := pkg.CreateEndpoints(svc)
	http.ListenAndServe(":80", pkg.NewHTTPHandler(endpoints))
}
