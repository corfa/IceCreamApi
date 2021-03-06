package realizationService

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

const (

	tokenTTL   = 2 * time.Hour
)
func (s *Service) CreateJWToken(id int) (string, error) {


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})

	return token.SignedString([]byte(os.Getenv("JWTKey")))
}

func (s *Service) ReadJWToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv("JWTKey")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}