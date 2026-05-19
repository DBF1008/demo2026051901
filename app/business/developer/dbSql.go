package developer

import (
	"fmt"
	"gofly/global"
	"gofly/utils/gf"
	"os"
)

func ExecSqlFile(tables []string, pathname string) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local&timeout=1000ms", global.App.Config.DBconf.Username, global.App.Config.DBconf.Password, global.App.Config.DBconf.Hostname, global.App.Config.DBconf.Hostport, global.App.Config.DBconf.Database)
	f, _ := os.Create(pathname)
	_ = gf.Dump(
		dsn,                   // DSN
		gf.WithDropTable(),
		gf.WithData(),         // Option: Dump Data (Default: Only dump table schema)
		gf.WithTables(tables), // Option: Dump Tables (Default: All tables)
		gf.WithWriter(f),      // Option: Writer (Default: os.Stdout)
	)
	f.Close()
}
