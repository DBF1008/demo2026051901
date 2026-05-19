package middleware

import (
	"fmt"
	"gofly/global"
	"gofly/utils/gf"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ValidityAPi() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := global.App.Config
		var apisecret = conf.App.Apisecret
		encrypt := c.Request.Header.Get("verify-encrypt")
		verifytime := c.Request.Header.Get("verify-time")
		mdsecret := gf.Md5(apisecret + verifytime)
		var NoVerifyAPIRoot_arr []string
		if global.App.Config.App.NoVerifyAPIRoot != "" {
			NoVerifyAPIRoot_arr = strings.Split(global.App.Config.App.NoVerifyAPIRoot, `,`)
		} else {
			NoVerifyAPIRoot_arr = make([]string, 0)
		}
		var NoVerifyAPI_arr []string
		if global.App.Config.App.NoVerifyAPI != "" {
			NoVerifyAPI_arr = strings.Split(global.App.Config.App.NoVerifyAPI, `,`)
		} else {
			NoVerifyAPI_arr = make([]string, 0)
		}
		rootPath := strings.Split(c.Request.URL.Path, "/")
		if (len(rootPath) > 2 && IsContain(NoVerifyAPIRoot_arr, rootPath[1])) || IsContain(NoVerifyAPI_arr, c.Request.URL.Path) || strings.Contains(c.Request.URL.Path, "/common/uploadfile/get_image") {
			c.Next()
		} else {
			verifytimeint, _ := strconv.ParseInt(verifytime, 10, 64)
			fmt.Println(time.Now().Unix(), verifytimeint)
			if mdsecret == encrypt && (time.Now().Unix()-verifytimeint < 60*15) {
				c.Next()
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code":    1,
					"message": "您的请求不合法，请按规范请求数据!",
					"result":  nil,
				})
				c.Abort()
			}
		}

	}
}
