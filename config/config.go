/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/3/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package config

import (
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"path/filepath"
)

var (
	Version = "1.0.0"

	LenStackBuf = 4096

	ModuleConfigRoot = ""
)

type BaseConfig struct {
	Cluster     ClusterConfig `yaml:"cluster"`
	PathLog     string        `yaml:"path.log"`
	PathProfile string        `yaml:"path.profile"`

	ConsolePort   int    `yaml:"console.port"`
	ConsolePrompt string `yaml:"console.prompt"`
}

type ClusterConfig struct {
	ID      string   `yaml:"node"`    //集群中的节点id 需要保证整个集群中唯一
	Name    string   `yaml:"name"`    //集群名称，不用业务使用不同的集群
	Servers []string `yaml:"servers"` //集群服务器列表
	Timeout uint     `yaml:"timeout"`

	Username string  `yaml:"username"`
	Password string  `yaml:"password"`

	TTL int64 `yaml:"ttl"` //

	//CertFile string
	//KeyFile  string
	//CommonName string
}

func (config ClusterConfig) IsValid() bool {
	return config.Name != "" && config.Servers != nil && len(config.Servers) != 0
}

func Init(configPath string) *BaseConfig {
	dir, _ := filepath.Abs(configPath)
	//log.Debugf("configuration path is %v", dir)
	ModuleConfigRoot = dir + string(filepath.Separator) + "modules" + string(filepath.Separator)

	config := &BaseConfig{}
	baseConfigPath := dir + string(filepath.Separator) + "aliensboot.yml"
	LoadConfigData(baseConfigPath, config)
	//fmt.Println("load config %v", config)
	return config
}

//func LoadConfigData(name string, config interface{}) {
//	if config == nil {
//		return
//	}
//	path := ModuleConfigRoot + name + ".yml"
//	data, err := ioutil.ReadFile(path)
//	if err != nil {
//		log.Fatalf("module %v config file is not found %v", name, path)
//	}
//	err = json.Unmarshal(data, config)
//	if err != nil {
//		log.Fatalf("load config %v err %v", path, err)
//	}
//}

func LoadConfigData(path string, config interface{}) {
	if config == nil {
		return
	}
	ymlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("config file is not found %v", path)
	}
	err = yaml.Unmarshal(ymlFile, config)
	if err != nil {
		log.Fatalf("load config %v err %v", path, err)
	}
	//log.Debugf("load config data %v - %v", path, config)
}

func LoadModuleConfigData(name string, config interface{}) {
	if config == nil {
		return
	}
	path := ModuleConfigRoot + name + ".yml"
	LoadConfigData(path, config)
}
