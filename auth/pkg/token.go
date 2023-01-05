package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var ErrTokenExpired = jwt.ErrTokenExpired

type Token interface {
	// Validates token
	IsValid() bool

	// Get token's associated user
	GetUser() *User
}

type TokenGenerator interface {
	// Generates a token
	Sign(user *User, exp time.Time, secret string) (string, error)

	// Validates token
	Verify(token, secret string) (Token, error)
}

// jwtToken implements Token
type jwtToken struct {
	token *jwt.Token
}

// Validates token issuer and audience
func (t *jwtToken) IsValid() bool {
	claims := t.token.Claims.(jwt.MapClaims)
	return claims.VerifyAudience("renting", true) && claims.VerifyIssuer("auth", true)
}

// Returns token payload
func (t *jwtToken) GetUser() *User {
	claims := t.token.Claims.(jwt.MapClaims)
	return claims["user"].(*User)
}

type jwtGenerator struct{}

func NewTokenGenerator() TokenGenerator {
	return &jwtGenerator{}
}

func (t *jwtGenerator) Sign(user *User, exp time.Time, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":  "renting",
		"sub":  user.ID,
		"exp":  exp.Unix(),
		"iss":  "auth",
		"user": user,
	})
	return token.SignedString([]byte(secret))
}

func (t *jwtGenerator) Verify(tokenStr, secret string) (Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	return &jwtToken{token}, err
}
