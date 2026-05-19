package main

import (
	"gofly/bootstrap"
	"gofly/global"
	"runtime"
	"strconv"
)

func main() {
	global.App.Config.InitializeConfig()
	global.App.Log = bootstrap.InitializeLog()
	global.App.Log.Info("项目启动成功")
	cpu_num, _ := strconv.Atoi(global.App.Config.App.CPUnum)
	mycpu := runtime.NumCPU()
	if cpu_num > mycpu {
		cpu_num = mycpu
	}
	if cpu_num > 0 {
		runtime.GOMAXPROCS(cpu_num)
	} else {
		runtime.GOMAXPROCS(mycpu)
	}

	bootstrap.RunServer()
}
