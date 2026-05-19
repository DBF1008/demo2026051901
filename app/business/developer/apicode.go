package developer

import (
	"encoding/json"
	"gofly/model"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type Apicode struct {
}

func init() {
	fpath := Apicode{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

func (api *Apicode) Installcode(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	data, err := model.DB().Table("common_apidoc").Where("id", parameter["id"]).Fields("cid,url,getdata_type,tablename,apicode_type,is_install,fields,method").First()
	if err != nil {
		results.Failed(c, "生成api接口代码失败", err)
	} else {
		type_id, _ := model.DB().Table("common_apidoc_group").Where("id", data["cid"]).Value("type_id")
		rooturl, _ := model.DB().Table("common_apidoc_type").Where("id", type_id).Value("model_name")
		model_name := "business"
		model_name_str := gf.InterfaceTostring(rooturl)
		if model_name_str != "" {
			model_name = model_name_str
		}
		if data["url"] == "" {
			results.Failed(c, "url地址为空", nil)
		} else if data["tablename"] == "" {
			results.Failed(c, "数据库表不能为空,选择数据表表提交保存再生成代码", nil)
		} else {
			CreatApicodeFile(model_name, data)
			model.DB().Table("common_apidoc").
				Data(map[string]interface{}{"is_install": 1}).
				Where("id", parameter["id"]).
				Update()
			results.Success(c, "生成api接口代码成功！", data, nil)
		}
	}
}

func (api *Apicode) Uninstallcode(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	data, err := model.DB().Table("common_apidoc").Where("id", parameter["id"]).Fields("cid,url,getdata_type,tablename,apicode_type,is_install,fields,method").First()
	if err != nil {
		results.Failed(c, "卸载失败", err)
	} else {
		UnApicodeFile(data)
		model.DB().Table("common_apidoc").Data(map[string]interface{}{"is_install": 2}).Where("id", parameter["id"]).Update()
		results.Success(c, "卸载成功！", data, nil)
	}
}

func (api *Apicode) RemoveFile(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	data, err := model.DB().Table("common_apidoc").Where("id", parameter["id"]).Fields("cid,url,getdata_type,tablename,apicode_type,is_install,fields,method").First()
	if err != nil {
		results.Failed(c, "删除文件失败", err)
	} else {
		type_id, _ := model.DB().Table("common_apidoc_group").Where("id", data["cid"]).Value("type_id")
		rooturl, _ := model.DB().Table("common_apidoc_type").Where("id", type_id).Value("model_name")
		model_name := "business"
		model_name_str := gf.InterfaceTostring(rooturl)
		if model_name_str != "" {
			model_name = model_name_str
		}
		url := data["url"].(string)
		url_arr := strings.Split(url, `/`)
		filename := url_arr[len(url_arr)-1]
		model_path := strings.Split(url, filename)
		haselist, _ := model.DB().Table("common_apidoc").Where("url", "like", model_path[0]+"%").Where("is_install", 1).Count("*")
		if haselist == 0 {
			RemoveModel(model_name, data)
		}
		model.DB().Table("common_apidoc").Data(map[string]interface{}{"is_install": 0}).Where("id", parameter["id"]).Update()
		results.Success(c, "删除文件成功！", data, haselist)
	}
}
