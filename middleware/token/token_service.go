package middleware

import (
	"errors"
	"filestore/config"
	"filestore/payload"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type MyClaims struct {
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return config.MySecret, nil
	}
}

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		myClaims, err := ParseToken(token)
		if err != nil {
			log.Errorf("token无效:%v", err)
			c.JSON(http.StatusForbidden, payload.FailPayload("token无效"))
			//c.Abort()
			return
		}
		c.Set("Phone", myClaims.Phone)
		c.Next()
		return
	}
}

//func CheckToken(c *gin.Context) (myClaims *MyClaims, err error) {
//	token := c.GetHeader("token")
//	myClaims, err = ParseToken(token)
//	if err != nil {
//		log.Errorf("token无效:%v", err)
//	}
//	return
//}

func ParseToken(tokens string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokens, &MyClaims{}, Secret())
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

// 生成token
func MakeToken(phone string) (tokenString string, err error) {
	claim := MyClaims{
		Phone: phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(4 * time.Hour * time.Duration(1))), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                       // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString(config.MySecret)
	return tokenString, err
}