package middleware

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	jwtMiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/lightserver-cn/golang-utils/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

const JwtKey = "ls"

func GenerateJwtMiddleware() *jwtMiddleware.Middleware {
	return jwtMiddleware.New(jwtMiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (i interface{}, err error) {
			return []byte(JwtKey), nil
		},
		ContextKey: "",
		ErrorHandler: func(ctx iris.Context, err error) {
			if _, err = ctx.JSON(response.JsonResponse{
				ErrorCode: -1,
				Message:   "认证失败，请重新登录认证",
				Data:      nil,
				Success:   false,
			}); err != nil {
				logrus.Error(err)
			}
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}

func ParseToken(tokenString string, key string) (interface{}, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		return err, false
	}
}

func GetToken(ctx iris.Context) string {
	token := ctx.GetHeader("Authorization")
	if token != "" && len(token) > 7 {
		token = token[7:]
	}
	return token
}

func GetUserID(token string) (id uint64) {
	if token != "" && token != "undefined" && len(token) > 7 {
		v, _ := ParseToken(token, JwtKey)
		if v != "" {
			id = uint64(cast.ToInt(v.(jwt.MapClaims)["id"]))
		}
	}
	return
}
