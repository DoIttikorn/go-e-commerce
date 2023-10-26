package auth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Doittikorn/go-e-commerce/config"
	"github.com/Doittikorn/go-e-commerce/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokeyType string

const (
	Access  TokeyType = "access"
	Refresh TokeyType = "refresh"
	Admin   TokeyType = "admin"
	ApiKey  TokeyType = "apiKey"
)

type auth struct {
	mapClaims *authMapClaims
	cfg       config.JWTConfigImpl
}

type admin struct {
	*auth
}

type authApiKey struct {
	*auth
}

type authMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type AuthImpl interface {
	SignToken() string
}

type AdminImpl interface {
	SignToken() string
}

// คำนวณเวลาที่จะหมดอายุ
func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

// คำนวณเวลาที่สร้าง token
func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

// สร้าง token
func New(tokenType TokeyType, cfg config.JWTConfigImpl, claims *users.UserClaims) (AuthImpl, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	case ApiKey:
		return newApiKey(cfg), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

// check token ว่าถูกสร้างขึ้นโดยเราหรือไม่
func ParseToken(cfg config.JWTConfigImpl, tokenString string) (*authMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token is malformed")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token is expired")
		} else {
			return nil, fmt.Errorf("parse token failed : %v", err)
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseAdminToken(cfg config.JWTConfigImpl, tokenString string) (*authMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token is malformed")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token is expired")
		} else {
			return nil, fmt.Errorf("parse token failed : %v", err)
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func RepeatToken(cfg config.JWTConfigImpl, claims *users.UserClaims, exp int64) string {
	obj := &auth{
		cfg: cfg,
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "go-e-commerce",
				Subject:   "refresh-token",
				Audience:  []string{"user", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func (a *auth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	signedToken, _ := token.SignedString([]byte(a.cfg.SecretKey()))
	return signedToken
}

func (a *admin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	signedToken, _ := token.SignedString([]byte(a.cfg.AdminKey()))
	return signedToken
}

func (a *authApiKey) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	signedToken, _ := token.SignedString([]byte(a.cfg.ApiKey()))
	return signedToken
}

func newAccessToken(cfg config.JWTConfigImpl, claims *users.UserClaims) AuthImpl {
	return &auth{
		cfg: cfg,
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "go-e-commerce",
				Subject:   "access-token",
				Audience:  []string{"user", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}

}

func newRefreshToken(cfg config.JWTConfigImpl, claims *users.UserClaims) AuthImpl {
	return &auth{
		cfg: cfg,
		mapClaims: &authMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "go-e-commerce",
				Subject:   "refresh-token",
				Audience:  []string{"user", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}

}

func newAdminToken(cfg config.JWTConfigImpl) AuthImpl {
	return &admin{
		auth: &auth{
			cfg: cfg,
			mapClaims: &authMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "go-e-commerce",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}

}

func ParseApiKey(cfg config.JWTConfigImpl, tokenString string) (*authMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.ApiKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*authMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func newApiKey(cfg config.JWTConfigImpl) AuthImpl {
	return &authApiKey{
		auth: &auth{
			cfg: cfg,
			mapClaims: &authMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "go-e-commerce",
					Subject:   "api-key",
					Audience:  []string{"admin", "customer"},
					ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(2, 0, 0)),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
