package config

import (
	"flag"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	V *viper.Viper
}

// 配置文件位置
const confPath = "../../configs/"

// 初始化viper配置
func Viperize() (v *viper.Viper, err error) {
	v = viper.New()
	if err = bindFlags(v); nil != err {
		return
	}

	if err = loadFile(v); nil != err {
		return
	}

	configureViper(v)
	return
}

func bindFlags(v *viper.Viper) error {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := v.BindPFlags(pflag.CommandLine); nil != err {
		return err
	}

	return nil
}

func loadFile(v *viper.Viper) error {
	env := v.GetString("env")
	if "" == env {
		env = "dev"
	}

	fileName := "config." + env
	v.SetConfigName(fileName)
	v.AddConfigPath(confPath)
	return v.ReadInConfig()
}

func configureViper(v *viper.Viper) {
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
}
