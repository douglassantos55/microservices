package pkg

import (
	"github.com/go-kit/log"
)

type logging struct {
	next   Service
	logger log.Logger
}

func LoggingMiddleware(svc Service, logger log.Logger) Service {
	return &logging{svc, logger}
}

func (l *logging) Login(user string, pass string) (res *AuthResponse, err error) {
	defer func() {
		l.logger.Log(
			"method", "Login",
			"user", user,
			"pass", pass,
			"res", res,
			"err", err,
		)
	}()
	return l.next.Login(user, pass)
}

func (l *logging) Verify(token string) (user *User, err error) {
	defer func() {
		l.logger.Log(
			"method", "Verify",
			"token", token,
			"user", user,
			"err", err,
		)
	}()
	return l.next.Verify(token)
}
