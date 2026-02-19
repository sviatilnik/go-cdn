package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sviatilnik/go-cdn/internal/user"
)

var ErrInvalidToken = errors.New("invalid token")

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

func (s *AuthService) VerifyAccessToken(accessToken string) (*user.User, error) {
	token, err := jwt.Parse(
		accessToken,
		s.keyFunc(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(s.issuer),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, fmt.Errorf("parse token failed: %w", err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	email, _ := claims["user_email"].(string)
	name, _ := claims["user_name"].(string)

	return &user.User{
		Email: email,
		Name:  name,
	}, nil
}

func (s *AuthService) keyFunc() jwt.Keyfunc {
	return func(_ *jwt.Token) (interface{}, error) { return []byte(s.secret), nil }
}
