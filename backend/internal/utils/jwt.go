package utils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type JwtUser struct {
	ID string `json:"id"`
}

type TokenPair struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func (j *Auth) generateTokenPair(user *JwtUser) (TokenPair, error) {
	// 1.Generate acces token
	accessToken, err := j.generateToken(user, j.TokenExpiry)
	if err != nil {
		return TokenPair{}, err
	}

	// 2. Generate refresh token
	refreshToken, err := j.generateToken(user, j.RefreshExpiry)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *Auth) generateToken(user *JwtUser, ttl time.Duration) (string, error) {
	// Creaet claims
	claims := Claims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{j.Audience},
			Issuer:    j.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(ttl)),
		},
	}

	// Create a token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create a signed token
	signedToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *Auth) GetRefreshCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     j.CookieName,
		Path:     j.CookiePath,
		Value:    refreshToken,
		Expires:  time.Now().UTC().Add(j.RefreshExpiry),
		MaxAge:   int(j.RefreshExpiry.Seconds()),
		Domain:   j.CookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
