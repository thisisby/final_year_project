package jwt

import (
	"errors"
	"time"
)

type Config struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

var defaultConfig *Config

func MustInitializeConfig(
	secretKey string,
	accessTokenTTL,
	refreshTokenTTL time.Duration,
) {
	defaultConfig = &Config{
		SecretKey:       secretKey,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
	}
}

func GetConfig() (*Config, error) {
	if defaultConfig == nil {
		return nil, errors.New("JWT config is not initialized")
	}
	return defaultConfig, nil
}
