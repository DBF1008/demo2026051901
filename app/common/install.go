package common

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"gofly/global"
	"gofly/model"
	"gofly/utils/gf"
	"gofly/utils/results"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Install struct {
}

func init() {
	fpath := Install{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

func (api *Install) Index(context *gin.Context) {
	path, err := os.Getwd()
	if err != nil {
		results.Failed(context, "项目路径获取失败", nil)
		return
	}
	filePath := filepath.Join(path, "/resource/developer/template/install.lock")
	if _, err := os.Stat(filePath); err == nil {
		context.HTML(http.StatusOK, "isinstall.html", gin.H{
			"title": "已经安装页面",
		})
	} else {
		context.HTML(http.StatusOK, "install.html", gin.H{
			"title": "安装页面",
		})
	}

}

func (api *Install) Save(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	var parameter map[string]interface{}
	_ = json.Unmarshal(body, &parameter)
	path, err := os.Getwd()
	if err != nil {
		results.Failed(c, "项目路径获取失败", nil)
		return
	}
	model.CreateDataBase(parameter["username"], parameter["password"], parameter["hostname"], parameter["hostport"], parameter["database"])
	cferr := gf.UpConfFieldData(path, parameter)
	if cferr != nil {
		results.Failed(c, "修改数据库配置失败", nil)
		return
	}
	model.MyInit(2)
	time.Sleep(time.Second * 3)
	SqlPath := filepath.Join(path, "/resource/developer/template/gofly_basedb.sql")
	sqls, sqlerr := os.ReadFile(SqlPath)
	if sqlerr != nil {
		results.Failed(c, "数据库文件不存在："+SqlPath, nil)
		return
	}
	sqlArr := strings.Split(string(sqls), ";")
	for _, sql := range sqlArr {
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		model.ExecSql(sql)
	}
	if parameter["isInstalladmin"] == "install" {
		adminSqlPath := filepath.Join(path, "/resource/developer/template/admin_db.sql")
		adminsqls, adminsqlerr := os.ReadFile(adminSqlPath)
		if adminsqlerr != nil {
			results.Failed(c, "数据库文件不存在："+adminSqlPath, nil)
			return
		}
		adminSqlArr := strings.Split(string(adminsqls), ";")
		for _, adminsql := range adminSqlArr {
			adminsql = strings.TrimSpace(adminsql)
			if adminsql == "" {
				continue
			}
			model.ExecSql(adminsql)
		}
	}
	salt := time.Now().Unix()
	businesspass := fmt.Sprintf("%v%v", gf.Md5(parameter["businessPassword"].(string)), salt)
	model.DB().Table("business_account").Data(map[string]interface{}{"username": parameter["businessUsername"], "password": gf.Md5(businesspass), "salt": salt}).Where("id", 1).Update()
	if parameter["isInstalladmin"] == "install" {
		adminpass := fmt.Sprintf("%v%v", gf.Md5(parameter["adminPassword"].(string)), salt)
		model.DB().Table("admin_account").Data(map[string]interface{}{"username": parameter["adminUsername"], "password": gf.Md5(adminpass), "salt": salt}).Where("id", 1).Update()
	}
	filePath := filepath.Join(path, "/resource/developer/template/install.lock")
	os.Create(filePath)
	if _, ok := parameter["vuepath"]; ok && parameter["vuepath"] != "" {
		parameter["vueobjroot"] = filepath.Join(gf.InterfaceTostring(parameter["vuepath"]), "/business")
		file_path := fmt.Sprintf("%v", parameter["vuepath"])
		if _, err := os.Stat(file_path); err != nil {
			if !os.IsExist(err) {
				os.MkdirAll(file_path, os.ModePerm)
			}
		}
		vuesoure_path := filepath.Join(path, "/resource/developer/template/vuecode/")
		CopyDir(vuesoure_path, file_path)
		business_vue_path := filepath.Join(file_path, "/business.zip")
		admin_vue_path := filepath.Join(file_path, "/admin.zip")
		Unzip(business_vue_path, file_path)
		if parameter["isInstalladmin"] == "install" {
			Unzip(admin_vue_path, file_path)
		} else {
			app_admin_path := filepath.Join(path, "/app/admin")
			os.RemoveAll(app_admin_path)
			ChecAdminRemoveController()
		}
		os.RemoveAll(business_vue_path)
		os.RemoveAll(admin_vue_path)
	}
	results.Success(c, "安装成功,去前端刷新试试！", parameter, nil)
}

func ChecAdminRemoveController() {
	filePath := filepath.Join("app/controller.go")
	con_path := "gofly/app/admin"
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	var result = ""
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		if strings.Contains(string(a), con_path) {
			continue
		} else {
			result += string(a) + "\n"
		}
	}
	fw, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	w := bufio.NewWriter(fw)
	w.WriteString(result)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func DeCompress(zipFile, dest string) (string, error) {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			global.App.Log.Info(fmt.Sprintf("解压文件关闭失败: %v\n", err.Error()))
		}
	}()
	var (
		first string
		order int    = 0
	)
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return "", err
		}
		filename := filepath.Join(dest, file.Name)
		if order == 0 {
			first = filename
		}
		order += 1
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(filename, 0755)
			if err != nil {
				return "", err
			}
		} else {
			w, err := os.Create(filename)
			if err != nil {
				return "", err
			}
			//defer w.Close()
			_, err = io.Copy(w, rc)
			if err != nil {
				return "", err
			}
			iErr := w.Close()
			if iErr != nil {
				global.App.Log.Info(fmt.Sprintf("[unzip]: close io %s\n", iErr.Error()))
			}
			fErr := rc.Close()
			if fErr != nil {
				global.App.Log.Info(fmt.Sprintf("[unzip]: close io %s\n", fErr.Error()))
			}
		}
	}
	return first, nil
}

func Unzip(zipath, dir string) error {
	// Open zip file.
	reader, err := zip.OpenReader(zipath)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if err := unzipFile(file, dir); err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(file *zip.File, dir string) error {
	// Prevent path traversal vulnerability.
	name := strings.TrimPrefix(filepath.Join(string(filepath.Separator), file.Name), string(filepath.Separator))
	filePath := path.Join(dir, name)

	// Create the directory of file.
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// Open the file.
	r, err := file.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	// Create the file.
	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	// Save the decompressed file content.
	_, err = io.Copy(w, r)
	return err
}

func CopyDir(targetPath string, destPath string) error {
	err := filepath.Walk(targetPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		destPath := filepath.Join(destPath, path[len(targetPath):])
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		return copyFile(path, destPath)

	})
	return err
}

func copyFile(srcFile, destFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()
	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}
	err = dest.Sync()
	if err != nil {
		return err
	}

	return nil
}
