package myjwt

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var (
	secretKey string
	exp       time.Duration
)

const issuer = "finalai"

type MyData struct {
	// UserID string `json:"user_id"`
	Username string `json:"username"`
}

type MyClaims struct {
	MyData
	jwt.RegisteredClaims
}

func Init() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file")
	}

	secretKey = os.Getenv("JWT_SECRET_KEY")
	if len(secretKey) == 0 {
		panic("JWT_SECRET_KEY is not set in .env file")
	}

	expStr := os.Getenv("JWT_EXPIRATION")
	if len(expStr) == 0 {
		panic("JWT_EXPIRATION is not set in .env file")
	}

	hours, err := strconv.Atoi(expStr)
	if err != nil {
		panic("JWT_EXPIRATION must be a valid integer representing hours")
	}
	exp = time.Duration(hours) * time.Hour
}

func GenerateToken(data MyData) (string, error) {
	now := time.Now()
	expAt := now.Add(exp)

	claims := MyClaims{
		MyData: data,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   data.Username,
			ExpiresAt: jwt.NewNumericDate(expAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ParseToken(tokenStr string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}
