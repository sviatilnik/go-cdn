package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sviatilnik/go-cdn/internal/user"
)

type AuthService struct {
	issuer string
	secret string
	exp    int
}

func NewAuthService(issuer string, secret string, exp int) *AuthService {
	return &AuthService{
		issuer: issuer,
		secret: secret,
		exp:    exp,
	}
}

func (s *AuthService) CreateAccessToken(u *user.User) (string, error) {
	claims := jwt.MapClaims{
		"iss":        s.issuer,
		"sub":        u.Email,
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Unix() + int64(s.exp),
		"user_email": u.Email,
		"user_name":  u.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secret))
}
