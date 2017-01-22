package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config 配置结构体
type Config struct {
	Version    string     `json:"version"`              // 版本号
	Tool       string     `json:"tool"`                 // 编译工具：go、gb...
	Importpath string     `json:"importpath,omitempty"` // 待替换的变量所在包导入路径，如github.com/voidint/gbb/build
	Variables  []Variable `json:"variables,omitempty"`  // 待替换变量集合
	Debug      bool       `json:"-"`
}

// Variable 变量结构体
type Variable struct {
	Variable string `json:"variable"`
	Value    string `json:"value"`
}

// Load 从磁盘文件加载配置信息
func Load(filename string) (conf *Config, err error) {
	var c Config
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

// Save 将配置信息序列化成JSON并写到指定磁盘位置
func Save(conf *Config, filename string) (err error) {
	b, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)
}
