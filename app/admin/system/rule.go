package system

import (
	"encoding/json"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"reflect"
	"time"

	"gofly/utils/gform"

	"github.com/gin-gonic/gin"
)

type Rule struct {
}

func init() {
	fpath := Rule{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

func (api *Rule) Get_list(c *gin.Context) {
	menuList, _ := model.DB().Table("admin_auth_rule").Order("orderNo asc").Get()
	if menuList == nil {
		menuList = make([]gform.Data, 0)
	}
	for _, val := range menuList {
		if val["title"] == "" {
			val["title"] = val["locale"]
		}
	}
	menuList = gf.GetRuleTreeArray(menuList, 0, "")
	results.Success(c, "获取全部菜单列表", menuList, nil)
}

func (api *Rule) Get_parent(c *gin.Context) {
	menuList, err := model.DB().Table("admin_auth_rule").WhereIn("type", []interface{}{0, 1}).Fields("id,pid,title,locale").Order("orderNo asc").Get()
	if err != nil {
		results.Failed(c, "获取选项列表失败", err)
	} else {
		if menuList == nil {
			menuList = make([]gform.Data, 0)
		}
		for _, val := range menuList {
			if val["title"] == "" {
				val["title"] = val["locale"]
			}
		}
		menuList = gf.GetMenuChildrenArray(menuList, 0, "pid")
		results.Success(c, "菜单父级数据！", menuList, nil)
	}
}

func (api *Rule) Save(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	parameter["uid"] = user.ID
	var f_id float64 = 0
	if parameter["id"] != nil {
		f_id = parameter["id"].(float64)
	}
	if f_id == 0 {
		parameter["createtime"] = time.Now().Unix()
		delete(parameter, "id")
		addId, err := model.DB().Table("admin_auth_rule").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(c, "添加菜单失败", err)
		} else {
			if addId != 0 {
				model.DB().Table("admin_auth_rule").
					Data(map[string]interface{}{"orderNo": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(c, "添加成功！", addId, nil)
		}
	} else {
		delete(parameter, "children")
		delete(parameter, "spacer")
		res, err := model.DB().Table("admin_auth_rule").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(c, "更新菜单失败", err)
		} else {
			results.Success(c, "更新成功！", res, nil)
		}
	}
}

func (api *Rule) UpStatus(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := model.DB().Table("admin_auth_rule").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
	if err != nil {
		results.Failed(c, "更新失败！", err)
	} else {
		msg := "更新成功！"
		if res2 == 0 {
			msg = "暂无数据更新"
		}
		results.Success(c, msg, res2, nil)
	}
}

func (api *Rule) Del(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := model.DB().Table("admin_auth_rule").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(c.Copy(), "删除菜单失败", err)
	} else {
		model.DB().Table("admin_auth_rule").WhereIn("pid", ids.([]interface{})).Delete()
		results.Success(c, "删除成功！", res2, nil)
	}
	c.Abort()
}
