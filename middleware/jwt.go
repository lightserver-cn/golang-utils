package middleware

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	jwtMiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/lightserver-cn/golang-utils/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// 过期时间，默认 10 小时
const expireTime = 10 * time.Hour * time.Duration(1)

type MyClaims struct {
	Uid       uint64 `json:"uid"`           // 用户 uid
	NikeName  string `json:"nik"`           // 用户昵称
	Audience  string `json:"aud,omitempty"` // 接收 jwt 的一方
	ExpiresAt int64  `json:"exp,omitempty"` // jwt 的过期时间，这个过期时间必须要大于签发时间
	Id        string `json:"jti,omitempty"` // jwt 的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。
	IssuedAt  int64  `json:"iat,omitempty"` // jwt 的签发时间
	Issuer    string `json:"iss,omitempty"` // jwt 签发者
	NotBefore int64  `json:"nbf,omitempty"` // 定义在什么时间之前，该 jwt 都是不可用的
	Subject   string `json:"sub,omitempty"` // jwt 所面向的用户
}

// GenerateToken 创建 token
func GenerateToken(jwtKey string, claims *MyClaims) string {
	// 设置过期时间，默认 10 小时
	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = time.Now().Add(expireTime).Unix()
	}
	// 设置签发时间
	if claims.IssuedAt == 0 {
		claims.IssuedAt = time.Now().Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": claims.Uid,
		"nik": claims.NikeName,
		"iss": claims.Issuer,
		"aud": claims.Audience,
		"exp": claims.ExpiresAt,
		"jti": claims.Id,
		"iat": claims.IssuedAt,
		"nbf": claims.NotBefore,
		"sub": claims.Subject,
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		logrus.Errorf("token SignedString failure :", err.Error())
	}

	return tokenString
}

// GenerateJwtMiddleware 验证 token
func GenerateJwtMiddleware(jwtKey string) *jwtMiddleware.Middleware {
	return jwtMiddleware.New(jwtMiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (i interface{}, err error) {
			return []byte(jwtKey), nil
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

// ParseToken 解析 token
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

// GetToken 获取 token 字段
func GetToken(ctx iris.Context) string {
	token := ctx.GetHeader("Authorization")
	if token != "" && len(token) > 7 {
		token = token[7:]
	}
	return token
}

// 获取登陆 uid
func GetUserID(token, jwtKey string) uint64 {
	var uid uint64
	if token != "" && token != "undefined" && len(token) > 7 {
		v, _ := ParseToken(token, jwtKey)
		if v != "" {
			return uint64(cast.ToInt(v.(jwt.MapClaims)["uid"]))
		}
	}
	return uid
}
