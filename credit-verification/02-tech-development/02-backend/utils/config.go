// utils/config.go
package utils

import (
	"log"

	"github.com/spf13/viper"
)

// GlobalConfig 全局配置结构体（与config.yaml对应）
var GlobalConfig struct {
	Server struct {
		Port string `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`
	MySQL struct {
		DSN             string `mapstructure:"dsn"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"mysql"`
	Ethereum struct {
		RpcUrl             string `mapstructure:"rpc_url"`
		CreditContractAddr string `mapstructure:"credit_contract_addr"`
		PrivateKey         string `mapstructure:"private_key"`
	} `mapstructure:"ethereum"`
	JWT struct {
		Secret      string `mapstructure:"secret"`
		ExpireHours int    `mapstructure:"expire_hours"`
	} `mapstructure:"jwt"`
}

// InitConfig 初始化配置（读取config.yaml）
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config") // 配置文件在config/目录下

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}
	log.Println("配置初始化完成")
}
