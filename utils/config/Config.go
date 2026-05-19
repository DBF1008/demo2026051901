package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	DBconf DBconf `yaml:"dbconf"`
	App    App    `yaml:"app"`
	Jwt    Jwt    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Log    Log    `mapstructure:"log" json:"log" yaml:"log"`
}

func (config *Config) InitializeConfig() *Config {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// fmt.Println("path=", path)
	vip := viper.New()
	vip.AddConfigPath(path + "/resource")
	vip.SetConfigName("config")
	vip.SetConfigType("yaml")
	if err := vip.ReadInConfig(); err != nil {
		panic(err)
	}
	vip.WatchConfig()
	vip.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		if err := vip.Unmarshal(&config); err != nil {
			fmt.Println(err)
		}
		// 注意：配置变更后需手动调用 model.MyInit() 重连数据库
		// 以及重新初始化Redis连接池
		fmt.Println("配置已重新加载，请手动重启服务使数据库连接生效")
	})

	err = vip.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	return config
}
func listenSignal() {
	go func() {
		cmd := exec.Command("gofly", "run", "daemon", "restart")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
		}
		defer stdout.Close()

		if err := cmd.Start(); err != nil {
			panic(err)
		}
		reader := bufio.NewReader(stdout)
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			fmt.Print(line)
		}

		if err := cmd.Wait(); err != nil {
			fmt.Println(err)
		}
		opBytes, _ := io.ReadAll(stdout)
		fmt.Print(string(opBytes))

	}()
}
