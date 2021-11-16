package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

type Config struct {
	Token    string `json:"Token" env:"API_TOKEN"`
	LogLevel string `json:"LogLevel" env:"LOG_LEVEL"`
	DbUrl    string `json:"DbUrl" env:"DB_URL"`
}

func NewConfig() (*Config, error) {
	bt, err := ioutil.ReadFile(".config")
	if err != nil {
		return nil, err
	}

	res := &Config{}
	return res, json.Unmarshal(bt, res)
}

func NewConfigEnv() (*Config, error) {
	res := &Config{}
	vals := reflect.ValueOf(res).Elem()
	fields := reflect.TypeOf(res).Elem()
	for i := 0; i < vals.NumField(); i++ {
		envName := fields.Field(i).Tag.Get("env")
		val, ok := os.LookupEnv(envName)
		if !ok {
			return nil, fmt.Errorf("failed to lookup env var %s", envName)
		}
		vals.Field(i).SetString(val)
	}
	return res, nil
}
