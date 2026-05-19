package middleware

import (
	"gofly/global"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserClaims struct {
	ID         int64  `json:"id"`
	Accountid  int64  `json:"accountid"`
	BusinessID int64  `json:"businessID"`
	Openid     string `json:"openid"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	jwt.StandardClaims
}

var (
	secret       = []byte("16849841325189456f489")
	refreshToken = []byte("gofly_refresh_secret_key_2024") // 刷新token专用密钥
)

var Expirre = "180"
var effectTime = time.Duration(getiInt()) * time.Minute
func getiInt() int64 {

	num := global.App.Config.App.TokenOutTime
	intnum, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 2 * 60
	} else {
		return intnum
	}
}

func TokenOutTime(claims *UserClaims) int64 {
	return time.Now().Add(effectTime).Unix()
}

func GenerateToken(claims *UserClaims) interface{} {
	claims.ExpiresAt = time.Now().Add(effectTime).Unix()
	sign, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		panic(err)
	}
	return sign
}

func JwtVerify(c *gin.Context) {
	var NoVerifyTokenRoot_arr []string
	if global.App.Config.App.NoVerifyTokenRoot != "" {
		NoVerifyTokenRoot_arr = strings.Split(global.App.Config.App.NoVerifyTokenRoot, `,`)
	} else {
		NoVerifyTokenRoot_arr = make([]string, 0)
	}
	var NoVerifyToken_arr []string
	if global.App.Config.App.NoVerifyToken != "" {
		NoVerifyToken_arr = strings.Split(global.App.Config.App.NoVerifyToken, `,`)
	} else {
		NoVerifyToken_arr = make([]string, 0)
	}
	rootPath := strings.Split(c.Request.URL.Path, "/")
	if len(rootPath) > 2 && IsContain(NoVerifyTokenRoot_arr, rootPath[1]) {
		return
	} else if IsContain(NoVerifyToken_arr, c.Request.URL.Path) {
		return
	}
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.GetHeader("authorization")
	}
	if token == "" {
		panic("token 不存在")
	}
	c.Set("user", ParseToken(token))
}

func ParseToken(tokenString string) *UserClaims {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		panic(err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		panic("The token is invalid")
	}
	return claims
}

func Refresh(tokenString string) interface{} {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return refreshToken, nil
	})
	if err != nil {
		panic(err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		panic("The token is invalid")
	}
	jwt.TimeFunc = time.Now
	claims.StandardClaims.ExpiresAt = time.Now().Add(effectTime).Unix()
	return GenerateToken(claims)
}
func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
