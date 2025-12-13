package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}
	jwtSecret = []byte(secret)
}

type Claims struct {
	UserId primitive.ObjectID `json:"user_id"`
	Email  string             `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userId primitive.ObjectID, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserId: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GetUserIdFromToken(tokenString string) (primitive.ObjectID, error) {
	claims, err := ValidateToken(tokenString)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return claims.UserId, nil
}
