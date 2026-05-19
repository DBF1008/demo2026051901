package system

import (
	"encoding/json"
	"fmt"
	"gofly/model"
	"gofly/route/middleware"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"reflect"
	"strings"
	"time"

	"gofly/utils/gform"

	"github.com/gin-gonic/gin"
)

type Role struct {
}

func init() {
	fpath := Role{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

func (api *Role) Get_list(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	status := c.DefaultQuery("status", "")
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	role_id, _ := model.DB().Table("admin_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	role_ids := gf.GetAllChilIds("admin_auth_role", role_id.([]interface{}))
	all_role_id := gf.MergeArr(role_id.([]interface{}), role_ids)
	MDB := model.DB().Table("admin_auth_role")
	roleList, _ := MDB.Where(func() {
		MDB.WhereIn("id", all_role_id)
		if name != "" {
			MDB.Where("name", "like", "%"+name+"%")
		}
		if status != "" {
			MDB.Where("status", status)
		}
	}).OrWhere(func() {
		MDB.Where("uid", user.ID)
		if name != "" {
			MDB.Where("name", "like", "%"+name+"%")
		}
		if status != "" {
			MDB.Where("status", status)
		}
	}).Order("weigh asc").Get()
	roleList = gf.GetTreeArray(roleList, 0, "")
	if roleList == nil {
		roleList = make([]gform.Data, 0)
	}
	results.Success(c, "获取拥有角色列表", roleList, all_role_id)
}

func (api *Role) Get_parent(c *gin.Context) {
	getuser, _ := c.Get("user")
	user := getuser.(*middleware.UserClaims)
	role_id, _ := model.DB().Table("admin_auth_role_access").Where("uid", user.ID).Pluck("role_id")
	role_ids := gf.GetAllChilIds("admin_auth_role", role_id.([]interface{}))
	all_role_id := gf.MergeArr(role_id.([]interface{}), role_ids)
	menuList, _ := model.DB().Table("admin_auth_role").WhereIn("id", all_role_id).OrWhere("uid", user.ID).Fields("id,pid,name").Order("weigh asc").Get()
	menuList = gf.GetMenuChildrenArray(menuList, 0, "pid")
	results.Success(c, "部门父级数据！", menuList, nil)
}

func (api *Role) Get_menuList(c *gin.Context) {
	id := c.DefaultQuery("id", "0")
	pid := c.DefaultQuery("pid", "0")
	MDB := model.DB().Table("admin_auth_rule").Where("status", 0)
	if id == "0" || pid == "0" {
		getuser, _ := c.Get("user")
		user := getuser.(*middleware.UserClaims)
		role_id, _ := model.DB().Table("admin_auth_role_access").Where("uid", user.ID).Pluck("role_id")
		menu_id, _ := model.DB().Table("admin_auth_role").WhereIn("id", role_id.([]interface{})).Pluck("rules")
		if !IsContain(menu_id.([]interface{}), "*") {
			getmenus := ArraymoreMerge(menu_id.([]interface{}))
			MDB = MDB.WhereIn("id", getmenus)
		}
	} else {
		pid, _ := model.DB().Table("admin_auth_role").Where("id", id).Value("pid")
		menu_id_str, _ := model.DB().Table("admin_auth_role").Where("id", pid).Value("rules")
		if !strings.Contains(menu_id_str.(string), "*") {
			getmenus := Axplode(menu_id_str)
			MDB = MDB.WhereIn("id", getmenus)
		}
	}
	menuList, _ := MDB.Fields("id,pid,title,locale").Order("orderNo asc").Get()
	for _, val := range menuList {
		if val["title"] == "" {
			val["title"] = val["locale"]
		}
		delete(val, "locale")
	}
	menuList = GetMenuChildrenArray(menuList, 0)
	results.Success(c, "获取菜单数据", menuList, nil)
}

func (api *Role) Save(c *gin.Context) {
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
	if parameter["menu"] != nil && parameter["menu"] != "*" {
		rules := gf.GetRulesID("admin_auth_role", "pid", parameter["menu"])
		rudata := rules.([]interface{})
		var rulesStr []string
		for _, v := range rudata {
			str := fmt.Sprintf("%v", v)
			rulesStr = append(rulesStr, str)
		}
		parameter["rules"] = strings.Join(rulesStr, ",")
		parameter["menu"] = gf.JSONMarshalToString(parameter["menu"])
	}
	parameter["createtime"] = time.Now().Unix()
	if f_id == 0 {
		delete(parameter, "id")
		addId, err := model.DB().Table("admin_auth_role").Data(parameter).InsertGetId()
		if err != nil {
			results.Failed(c, "添加失败", err)
		} else {
			if addId != 0 {
				model.DB().Table("admin_auth_role").
					Data(map[string]interface{}{"weigh": addId}).
					Where("id", addId).
					Update()
			}
			results.Success(c, "添加成功！", addId, nil)
		}
	} else {
		delete(parameter, "children")
		delete(parameter, "spacer")
		res, err := model.DB().Table("admin_auth_role").
			Data(parameter).
			Where("id", f_id).
			Update()
		if err != nil {
			results.Failed(c, "更新失败", err)
		} else {
			results.Success(c, "更新成功！", res, nil)
		}
	}
}

func (api *Role) UpStatus(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	res2, err := model.DB().Table("admin_auth_role").Where("id", parameter["id"]).Data(map[string]interface{}{"status": parameter["status"]}).Update()
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

func (api *Role) Del(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	ids := parameter["ids"]
	res2, err := model.DB().Table("admin_auth_role").WhereIn("id", ids.([]interface{})).Delete()
	if err != nil {
		results.Failed(c, "删除失败", err)
	} else {
		results.Success(c, "删除成功！", res2, nil)
	}
}
