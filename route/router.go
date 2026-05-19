package route

import (

	// _ "gofly/app/admin"
	// _ "gofly/app/business"
	// _ "gofly/app/wxoffi"
	// _ "gofly/app/common"
	// _ "gofly/app/wxapp"
	// _ "gofly/app/home"
	_ "gofly/app"
	"net/http"

	"gofly/global"
	"gofly/route/middleware"
	"strings"
	"time"

	"gofly/utils/gf"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	R := gin.Default()
	R.SetTrustedProxies([]string{"127.0.0.1"})
	R.Static("/resource", "./resource")
	R.Static("/webadmin", "./resource/webadmin")
	R.Static("/webbusiness", "./resource/webbusiness")
	R.LoadHTMLFiles("./resource/developer/template/install.html", "./resource/developer/template/isinstall.html")
	R.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, global.App.Config.App.Rootview)
	})
	gin.SetMode(global.App.Config.App.RunlogType)
	R.MaxMultipartMemory = 8 << 20 // 8 MiB
	var str_arr []string
	if global.App.Config.App.Allowurl != "" {
		str_arr = strings.Split(global.App.Config.App.Allowurl, `,`)
	} else {
		str_arr = []string{"http://localhost:8080"}
	}

	R.Use(cors.New(cors.Config{
		AllowOrigins: str_arr,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"X-Requested-With", "Content-Type", "Authorization", "Businessid", "verify-encrypt", "ignoreCancelToken", "verify-time"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	R.Use(gin.Logger(), middleware.CustomRecovery())
	R.Use(middleware.LimitHandler())
	R.Use(middleware.ValidityAPi())
	R.Use(middleware.JwtVerify)
	R.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		c.JSON(404, gin.H{"code": 404, "message": "您" + method + "请求地址：" + path + "不存在！"})
	})
	gf.Bind(R)
	return R
}
