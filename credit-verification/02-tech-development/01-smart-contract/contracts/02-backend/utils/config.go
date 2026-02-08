package utils

import (
	"log"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	MySQL    MySQLConfig    `mapstructure:"mysql"`
	Ethereum EthereumConfig `mapstructure:"ethereum"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MySQLConfig struct {
	Dsn             string `mapstructure:"dsn"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type EthereumConfig struct {
	RpcUrl             string `mapstructure:"rpc_url"`
	CreditContractAddr string `mapstructure:"credit_contract_addr"`
	RoleContractAddr   string `mapstructure:"role_contract_addr"`
	PrivateKey         string `mapstructure:"private_key"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

var GlobalConfig Config

// InitConfig 初始化配置
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config") // 配置文件路径

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
}
