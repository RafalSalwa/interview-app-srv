package auth

import (
	simpleHandler "github.com/RafalSalwa/interview-app-srv/pkg/simple_handler"
	"net/http"
	"time"
)

type Auth struct {
	AuthMethod  string    `mapstructure:"method"`
	APIKey      string    `mapstructure:"apiKey"`
	BearerToken string    `mapstructure:"bearerToken"`
	BasicAuth   BasicAuth `mapstructure:"basic"`
	JWTToken    JWTConfig `mapstructure:"jwt"`
}

type BasicAuth struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
type JWTConfig struct {
	Access  Token `mapstructure:"accessToken"`
	Refresh Token `mapstructure:"refreshToken"`
}
type Token struct {
	PrivateKey string        `mapstructure:"privateKey"`
	PublicKey  string        `mapstructure:"publicKey"`
	ExpiresIn  time.Duration `mapstructure:"expiresIn"`
	MaxAge     int           `mapstructure:"maxAge"`
}

func Authorization(h simpleHandler.HandlerFunc) http.HandlerFunc {
	at, _ := NewAuthMethod(h, "basic")
	return at.middleware(h)
}
