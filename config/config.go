package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Version   string     `json:"version"`             // 版本号
	Tool      string     `json:"tool"`                // 编译工具：go、gb...
	Package   string     `json:"package,omitempty"`   // 待替换的变量所在包名
	Variables []Variable `json:"variables,omitempty"` // 待替换变量集合
}

type Variable struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
}

func Load(filename string) (conf *Config, err error) {
	var c Config
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &c, json.Unmarshal(b, &c)
}

func Save(conf *Config, filename string) (err error) {
	b, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}
