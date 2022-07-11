package service

import (
	"fmt"
	"github.com/Celtech/ACME/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTService is an interface around our JWT auth package
type JWTService interface {
	GenerateToken(email string, isUser bool) string
	ValidateToken(token string) (*jwt.Token, error)
}
type authCustomClaims struct {
	Name string `json:"name"`
	User bool   `json:"user"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issure    string
}

// JWTAuthService constructor for the service
func JWTAuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issure:    "RykeLabs",
	}
}

func getSecretKey() string {
	secret := config.GetConfig().GetString("secret")
	if secret == "" {
		secret = "correct-horse-battery-staple" // Just for you Keith :)
	}
	return secret
}

// GenerateToken returns a JWT token from an email password combo
func (service *jwtServices) GenerateToken(email string, isUser bool) string {
	tokenTTL := config.GetConfig().GetInt("services.jwt.tokenTTL")
	if tokenTTL == 0 {
		tokenTTL = 30
	}

	claims := &authCustomClaims{
		email,
		isUser,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(tokenTTL)).Unix(),
			Issuer:    service.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

// ValidateToken verifies a JWT token authenticity
func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token %v", token.Header["alg"])

		}
		return []byte(service.secretKey), nil
	})
}
