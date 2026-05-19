package user

import (
	"encoding/json"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

func init() {
	gf.Register(&Index{}, reflect.TypeOf(Index{}).PkgPath())
}

type Index struct {
}

func (api *Index) Login(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["username"] == nil || parameter["password"] == nil {
		results.Failed(c, "请提交用户账号或密码！", nil)
		return
	}
	username := parameter["username"].(string)
	password := parameter["password"].(string)
	res, err := model.DB().Table("business_account").Fields("id,accountID,businessID,password,salt,name").Where("username", username).OrWhere("email", username).First()
	if res == nil || err != nil {
		results.Failed(c, "账号不存在！", nil)
		return
	}
	pass := gf.Md5(password + res["salt"].(string))
	if pass != res["password"] {
		results.Failed(c, "您输入的密码不正确！", pass)
		return
	}
	//token
	token := middleware.GenerateToken(&middleware.UserClaims{
		ID:             res["id"].(int64),
		Accountid:      res["accountID"].(int64),
		BusinessID:     res["businessID"].(int64),
		StandardClaims: jwt.StandardClaims{},
	})
	model.DB().Table("business_account").Where("id", res["id"]).Data(map[string]interface{}{"loginstatus": 1, "lastLoginTime": time.Now().Unix(), "lastLoginIp": gf.GetIp(c)}).Update()
	model.DB().Table("login_logs").
		Data(map[string]interface{}{"type": 1, "uid": res["id"], "out_in": "in",
			"createtime": time.Now().Unix(), "loginIP": gf.GetIp(c)}).Insert()
	results.Success(c, "登录成功返回token！", token, nil)
}

func (api *Index) RegisterUser(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["username"] == nil || parameter["password"] == nil {
		results.Failed(c, "请提交用户账号或密码！", nil)
		return
	}
	password := parameter["password"].(string)
	userdata, _ := model.DB().Table("business_account").Fields("id").Where("username", parameter["username"]).First()
	if userdata != nil {
		results.Failed(c, "账号已存在！", nil)
		return
	}
	userdata2, _ := model.DB().Table("business_account").Fields("id").Where("email", parameter["email"]).First()
	if userdata2 != nil {
		results.Failed(c, "邮箱已存在！", nil)
		return
	}
	rnd := rand.New(rand.NewSource(6))
	salt := strconv.Itoa(rnd.Int())
	pass := gf.Md5(password + salt)
	userid, err := model.DB().Table("business_account").Data(map[string]interface{}{
		"salt":       salt,
		"username":   parameter["username"],
		"password":   pass,
		"email":      parameter["email"],
		"avatar":     "resource/staticfile/avatar.png",
		"createtime": time.Now().Unix(),
	}).InsertGetId()
	if err != nil {
		results.Failed(c, "注册失败", err)
	} else {
		results.Success(c, "注册成功", userid, nil)
	}
}

func (api *Index) Get_userinfo(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	userdata, err := model.DB().Table("business_account").Fields("id,businessID,username,name,nickname,city,company,avatar,status").Where("id", user.ID).First()
	if err != nil {
		results.Failed(c, "查找用户数据！", err)
	} else {
		rooturl, _ := model.DB().Table("common_config").Where("keyname", "rooturl").Value("keyvalue")
		if userdata["avatar"] == nil {
			userdata["avatar"] = rooturl.(string) + "resource/staticfile/avatar.png"
		} else if !strings.Contains(userdata["avatar"].(string), "http") && rooturl != nil {
			userdata["avatar"] = rooturl.(string) + userdata["avatar"].(string)
		}
		results.Success(c, "获取用户信息", map[string]interface{}{
			"userId":       userdata["id"],
			"businessID":   userdata["businessID"],
			"username":     userdata["username"],
			"name":         userdata["name"],
			"avatar":       userdata["avatar"],
			"introduction": userdata["remark"],
			"nickname":     userdata["nickname"],
			"city":         userdata["city"],
			"company":      userdata["company"],
			"rooturl":      rooturl,
			"role":         "admin",
		}, nil)
	}
}

func (api *Index) Refreshtoken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	newtoken := middleware.Refresh(token)
	results.Success(c, "刷新token", newtoken, nil)
}

func (api *Index) Logout(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token != "" {
		middleware.Refresh(token)
		getuser, _ := c.Get("user")
		if getuser != nil {
			user := getuser.(*middleware.UserClaims)
			model.DB().Table("business_account").Where("id", user.ID).Data(map[string]interface{}{"loginstatus": 0}).Update()
		}
	}
	results.Success(c, "退出登录", true, nil)
}

func (api *Index) Get_code(c *gin.Context) {
	email := c.DefaultQuery("email", "")
	if email == "" {
		results.Failed(c, "请填写邮箱", nil)
	} else {
		emailConfig, _ := model.DB().Table("common_email").Where("data_from", "common").First()
		if emailConfig == nil {
			results.Failed(c, "请到admin后台“配置管理”配置邮箱", nil)
		} else {
			code := gf.GenValidateCode(6)
			sender := emailConfig["sender_email"].(string)
			authCode := emailConfig["auth_code"].(string)
			mailTitle := emailConfig["mail_title"].(string)
			mailBody := emailConfig["mail_body"].(string)

			m := gomail.NewMessage()
			m.SetHeader("From", sender)
			m.SetHeader("To", email)
			m.SetHeader("Subject", mailTitle)
			m.SetBody("text/html", mailBody)

			// zipPath := "./user/temp.zip"
			// m.Attach(zipPath)

			service_host := "smtp.qq.com"
			if _, ok := emailConfig["service_host"]; ok {
				service_host = emailConfig["service_host"].(string)
			}
			service_port := 587
			if _, ok := emailConfig["service_port"]; ok {
				service_port = gf.InterfaceToInt(emailConfig["service_port"])
			}
			d := gomail.NewDialer(service_host, service_port, sender, authCode)
			err := d.DialAndSend(m)
			_, erro := model.DB().Table("common_verify_code").Data(map[string]interface{}{
				"keyname":    email,
				"code":       code,
				"createtime": time.Now().Unix(),
			}).Insert()
			results.Success(c, "获取验证码", err, erro)
		}
	}
}

func (api *Index) ResetPassword(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	if parameter["code"] == nil || parameter["password"] == nil {
		results.Failed(c, "请提交验证码或密码！", nil)
		return
	}
	password := parameter["password"].(string)
	userdata2, _ := model.DB().Table("business_account").Where("email", parameter["email"]).Fields("id").First()
	if userdata2 == nil {
		results.Failed(c, "邮箱不存在！", nil)
		return
	}
	code, _ := model.DB().Table("common_verify_code").Where("keyname", parameter["email"]).Value("code")
	if code == nil || code != parameter["code"] {
		results.Failed(c, "验证码无效", nil)
		return
	}
	rnd := rand.New(rand.NewSource(6))
	salt := strconv.Itoa(rnd.Int())
	pass := gf.Md5(password + salt)
	res, err := model.DB().Table("business_account").Where("id", userdata2["id"]).Data(map[string]interface{}{
		"salt":     salt,
		"password": pass,
	}).Update()
	if err != nil {
		results.Failed(c, "重置密码失败", err)
	} else {
		results.Success(c, "重置密码成功", res, nil)
	}
}

func (api *Index) Get_logininfo(c *gin.Context) {
	res2, err := model.DB().Table("common_logininfo").Where("type", "business").OrWhere("type", "common").Fields("title,des,image").Order("weigh asc,id desc").Get()
	if err != nil {
		results.Failed(c, "获取登录页面失败", err)
	} else {
		results.Success(c, "获取登录页面成功！", res2, nil)
	}
}
