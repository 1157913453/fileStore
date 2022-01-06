package token_service

import (
	"errors"
	"filestore/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"time"
)

type MyClaims struct {
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}

var MySecret = []byte("天涯")

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte("天涯"), nil
	}
}

func ParseToken(tokenss string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenss, &MyClaims{}, Secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}

// MakeToken 生成token
func MakeToken(phone string) (tokenString string, err error) {
	claim := MyClaims{
		Phone: phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour * time.Duration(1))), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                       // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString(MySecret)
	return tokenString, err
}

func UpdateToken(phone, token string) error {
	return models.UpdateToken(phone, token)
}

func CheckToken(c *gin.Context) (myClaims *MyClaims, err error) {
	token := c.GetHeader("token")
	myClaims, err = ParseToken(token)
	if err != nil {
		log.Errorf("token无效:%v", err)
	}
	return
}
