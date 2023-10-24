package jwtauth

import (
	"time"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

// defaultExpired 默认过期时间
const defaultExpired = time.Hour * 24

var (
	// ErrSecretInvalid .
	ErrSecretInvalid = errors.New("the secret not be empty.")
	// ErrTokenExpired .
	ErrTokenExpired = errors.New("token has expired.")
	// ErrTokenInvalid .
	ErrTokenInvalid = errors.New("invalid token")
)

// UserClaims .
type UserClaims struct {
	ID string
}

// jwtStandardClaims .
type jwtStandardClaims struct {
	*UserClaims
	*jwt.StandardClaims
}

var opts = new(options)

// options .
type options struct {
	// secret signingKey
	secret []byte
	// expire token 过期时间
	expire time.Duration
}

type option func(o *options)

// Expire .
func Expire(d int) option {
	return func(o *options) {
		if d != 0 {
			o.expire = time.Second * time.Duration(d)
		}
	}
}

// Set .
func Set(secret string, os ...option) {
	if len(secret) == 0 {
		panic(ErrSecretInvalid)
	}

	var options = &options{secret: []byte(secret), expire: defaultExpired}
	for _, opt := range os {
		opt(options)
	}

	opts = options
}

// Decode .
func Decode(user *UserClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtStandardClaims{
		UserClaims: user,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(opts.expire).Unix(),
		},
	}).SignedString(opts.secret)
}

// Encode .
func Encode(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, new(jwtStandardClaims), func(token *jwt.Token) (interface{}, error) {
		return opts.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwtStandardClaims); ok && token.Valid {
		if time.Now().Unix() > claims.ExpiresAt {
			return nil, errors.New("token has expired")
		}
		return claims.UserClaims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
