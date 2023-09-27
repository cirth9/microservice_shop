package middleware

import (
	"Shop-api/user_web/global"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"log"
	"time"
)

// UserBasicClaims
// jwt claims
type UserBasicClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

// GetToken 获取token
func GetToken(nickname string, id, authorityId uint) (string, error) {
	expirationTime := time.Now().Add(global.TokenExpireDuration)
	log.Println("当前时间，token过期时间>>>>>>>>>>>>>", time.Now().Unix(), expirationTime.Unix())
	claims := &UserBasicClaims{
		ID:          id,
		NickName:    nickname,
		AuthorityId: authorityId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(global.Secret)
}

// ParseToken 解析token
func ParseToken(tokenString string) (*UserBasicClaims, error) {
	zap.S().Info("tokenString  ", tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &UserBasicClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return global.Secret, nil
	})
	if err != nil {
		zap.S().Info("error ", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if token != nil {
		zap.S().Info("token.claims", token.Claims)
		claims, ok := token.Claims.(*UserBasicClaims)
		if ok && token.Valid {
			zap.S().Info("claims ", claims)
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}
}
