package article

import (
	"gofly/utils/gf"
	"gofly/utils/results"
	"reflect"

	"github.com/gin-gonic/gin"
)

type wxapi struct {
}

func init() {
	fpath := wxapi{}
	gf.Register(&fpath, reflect.TypeOf(fpath).PkgPath())
}

func (api *wxapi) Get_data(c *gin.Context) {
	results.Success(c, "测试获取数据接口", "张三的数据", "扩展数据")
}
