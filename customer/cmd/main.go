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

	svc := pkg.NewService(pkg.NewValidator())
	svc = pkg.LoggingMiddleware(svc, logger)

	httpHandler := pkg.MakeHTTPHandler(svc, os.Getenv("AUTH_SERVICE_URL"))
	http.ListenAndServe(":80", httpHandler)
}
