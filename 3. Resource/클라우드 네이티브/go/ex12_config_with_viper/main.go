package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func localFileConfig() {
	viper.SetConfigFile("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/service/")
	viper.AddConfigPath("$HOME/.service")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error reading config: %w", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed:", in.Name)
	})
}

func remoteFileConfig() {
	viper.AddRemoteProvider("etcd", "http://127.0.0.1:4001", "/config/service.json")
	viper.SetConfigType("json")
	viper.ReadRemoteConfig()
}

func main() {

}
