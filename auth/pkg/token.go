package pkg

import "github.com/golang-jwt/jwt/v4"

const (
	JWT_SIGN_SECRET_ENV    = "JWT_SIGN_SECRET"
	JWT_REFRESH_SECRET_ENV = "JWT_REFRESH_SECRET"
)

var ErrTokenExpired = jwt.ErrTokenExpired

// Token payload
type Payload map[string]any

type TokenGenerator interface {
	// Generates a token with claims
	Sign(payload Payload, secret string) (string, error)

	// Validates token and returns claims
	Verify(token, secret string) (Payload, error)
}

type jwtGenerator struct{}

func NewTokenGenerator() TokenGenerator {
	return &jwtGenerator{}
}

func (t *jwtGenerator) Sign(payload Payload, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payload))
	return token.SignedString([]byte(secret))
}

func (t *jwtGenerator) Verify(tokenStr, secret string) (Payload, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	return Payload(token.Claims.(jwt.MapClaims)), err
}
