package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"reconcip.com.br/microservices/customer/pkg"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	repository, err := pkg.NewMongoRepository(
		os.Getenv("MONGODB_URL"),
		os.Getenv("MONGODB_USER"),
		os.Getenv("MONGODB_PASSWORD"),
	)

	if err != nil {
		panic(err)
	}

	svc := pkg.NewService(pkg.NewValidator(), repository)
	svc = pkg.NewLoggingService(svc, logger)

	httpHandler := pkg.MakeHTTPHandler(svc, os.Getenv("AUTH_SERVICE_URL"))
	http.ListenAndServe(":80", httpHandler)
}
