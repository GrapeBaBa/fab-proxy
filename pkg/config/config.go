package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var Conf *Config
var once sync.Once

type Node struct {
	Addr             string   `yaml:"addr" mapstructure:"addr"`
	Key              string   `yaml:"tls_key" mapstructure:"tls_key"`
	Cert             string   `yaml:"tls_cert" mapstructure:"tls_cert"`
	OverrideHostname string   `yaml:"override_hostname" mapstructure:"override_hostname"`
	RootCerts        []string `yaml:"root_certs" mapstructure:"root_certs"`
}

type Config struct {
	Crypto      Crypto `yaml:"crypto"`
	ServerType  string `yaml:"server_type" mapstructure:"server_type"`
	Port        string `yaml:"port" mapstructure:"port"`
	Peers       []Node `yaml:"peers" mapstructure:"peers"`
	Orderers    []Node `yaml:"orderers" mapstructure:"orderers"`
	Channel     string `yaml:"channel" mapstructure:"channel"`
	Concurrency int    `yaml:"concurrency" mapstructure:"concurrency"`
	AccountNum  int    `yaml:"account_num" mapstructure:"account_num"`
}

type Crypto struct {
	MSPID    string `yaml:"msp_id" mapstructure:"msp_id"`
	PrivKey  string `yaml:"priv_key" mapstructure:"priv_key"`
	SignCert string `yaml:"sign_cert" mapstructure:"sign_cert"`
	//TLSCACerts []string
}

func InitConfig(path string, filename string) {
	once.Do(func() {
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.SetConfigName(filename)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(path)
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}
		c := &Config{}
		err = viper.Unmarshal(c)
		if err != nil {
			panic("load config failed")
		}
		Conf = c
	})
}
