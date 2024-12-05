package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

type Claims struct {
	UserID  int
	IsAdmin bool
	jwt.RegisteredClaims
}

func GenerateTokenPair(userID int, isAdmin bool) (accessToken, refreshToken string, err error) {
	accessToken, err = GenerateToken(userID, isAdmin, AccessToken)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = GenerateToken(userID, isAdmin, RefreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func GenerateToken(userID int, isAdmin bool, tokenType string) (string, error) {
	config, err := GetConfig()
	if err != nil {
		return "", fmt.Errorf("jwt - GenerateToken - GetConfig: %w", err)
	}

	var claims Claims
	if tokenType == RefreshToken {
		claims = Claims{
			UserID:  userID,
			IsAdmin: isAdmin,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.AccessTokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
	} else {
		claims = Claims{
			UserID:  userID,
			IsAdmin: isAdmin,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshTokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

func ParseToken(tokenString string) (*Claims, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, fmt.Errorf("jwt - ParseToken - GetConfig: %w", err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt - ParseToken - token.Method %v", token.Header["alg"])
		}
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrTokenNotValidYet
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		default:
			return nil, ErrTokenInvalid
		}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}
