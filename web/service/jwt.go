package service

import (
	"fmt"
	"github.com/Celtech/ACME/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTService is an interface around our JWT auth package
type JWTService interface {
	GenerateToken(email string, isUser bool) string
	ValidateToken(token string) (bool, error)
}

type jwtServices struct {
	secretKey string
	issure    string
}

// JWTAuthService constructor for the service
func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issure:    getIssuer(),
	}
}

func getSecretKey() string {
	secret := config.GetConfig().GetString("secret")
	if secret == "" {
		secret = "correct-horse-battery-staple" // Just for you Keith :)
	}
	return secret
}

func getIssuer() string {
	issuer := config.GetConfig().GetString("services.jwt.issuer")
	if len(issuer) == 0 {
		issuer = "RykeLabs"
	}
	return issuer
}

// GenerateToken returns a JWT token from an email password combo
func (service *jwtServices) GenerateToken(email string, isUser bool) string {
	tokenTTL := config.GetConfig().GetInt("services.jwt.tokenTTL")
	if tokenTTL == 0 {
		tokenTTL = 30
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":    time.Now().Add(time.Minute * time.Duration(tokenTTL)).Unix(),
		"iss":    service.issure,
		"iat":    time.Now().Unix(),
		"email":  email,
		"isUser": isUser,
	})

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

// ValidateToken verifies a JWT token authenticity
func (service *jwtServices) ValidateToken(encodedToken string) (bool, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token %v", token.Header["alg"])
		}

		return []byte(service.secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["iss"] == getIssuer() {
		return true, err
	}

	return false, err
}
