package main

import (
	"net/http"
	"os"

	"github.com/go-kit/log"
	"reconcip.com.br/microservices/auth/pkg"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestamp)
	logger = log.WithPrefix(logger, "caller", log.DefaultCaller)

	svc := pkg.NewService(pkg.NewTokenGenerator())
	svc = pkg.NewLoggingService(svc, logger)

	http.ListenAndServe(":80", pkg.NewHTTPHandler(svc))
}
