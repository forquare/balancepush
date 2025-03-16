package config

import (
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"runtime"
	"sync"
)

var (
	once   sync.Once
	config Config
)

type Config struct {
	GoCardless    GoCardlessConfig    `yaml:"gocardless"`
	Pushover      PushoverConfig      `yaml:"pushover"`
	Requisitioner RequisitionerConfig `yaml:"requisitioner"`
}

type GoCardlessConfig struct {
	Credentials CredentialsConfig `yaml:"credentials"`
	Bank        BankConfig        `yaml:"bank"`
}

type CredentialsConfig struct {
	SecretID  string `yaml:"secretID"`
	SecretKey string `yaml:"secretKey"`
}

type BankConfig struct {
	Institution string        `yaml:"institution"`
	Accounts    []BankAccount `yaml:"accounts"`
}

type BankAccount struct {
	ID             string  `yaml:"id"`
	Name           string  `yaml:"name"`
	BalanceType    string  `yaml:"balanceType"`
	Balance        float64 `yaml:"balance,omitempty"`
	Currency       string  `yaml:"currency,omitempty"`
	CurrencySymbol string  `yaml:"currencySymbol,omitempty"`
}

type PushoverConfig struct {
	Tokens PushoverTokens `yaml:"tokens"`
}

type PushoverTokens struct {
	User string `yaml:"user"`
	App  string `yaml:"app"`
}

type RequisitionerConfig struct {
	Listen   ListenConfig   `yaml:"listen"`
	Redirect RedirectConfig `yaml:"redirect"`
}

type ListenConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RedirectConfig struct {
	Proto string `yaml:"proto"`
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Path  string `yaml:"path"`
}

func GetConfig() Config {
	once.Do(func() {
		logger.Debug("Reading config file")
		viper.SetConfigName("balancepush")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("../balancepush_bank_checker/")
		viper.AddConfigPath("$HOME/config/")
		viper.AddConfigPath("/usr/local/etc/")
		viper.AddConfigPath("/etc/")
		err := viper.ReadInConfig()
		if err != nil {
			logger.Println(err)
			logger.Fatal("Cannot read config file. File may not exist, or be in the wrong format.")
		}
		err = viper.Unmarshal(&config)
		if err != nil {
			logger.Fatal("Cannot read config file. File may be in the wrong format.")
		}
	})

	if logger.GetLevel() == logger.TraceLevel {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			logger.Tracef("Returning config to %s", details.Name())
		}
	}

	return config
}
