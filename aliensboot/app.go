/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/8/21
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package aliensboot

import (
	"flag"
	"fmt"
	"github.com/KylinHe/aliensboot-core/cluster/center"
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/console"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/module"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
)

var (
	debug      = false
	configPath = "" //配置文件根目录，默认当前
	tag        = ""
)

func init() {
	flag.BoolVar(&debug, "debug", false, "debug flag")
	flag.StringVar(&configPath, "config", "config", "configuration path")
	flag.StringVar(&tag, "tag", "aliensboot", "log tag")
	flag.Parse()

}

func IsDebug() bool {
	return debug
}

func Run(mods ...module.Module) {
	baseConfig := config.Init(configPath)

	//log.Debugf("config data %+v", baseConfig)
	log.Init(debug, tag, baseConfig.PathLog)
	clusterName := os.Getenv("ClusterName")
	if clusterName != "" {
		baseConfig.Cluster.Name = clusterName
	}
	clusterID := os.Getenv("ClusterNode")
	if clusterID != "" {
		baseConfig.Cluster.ID = clusterID
	}
	clusterAddress := os.Getenv("ClusterAddress")
	if clusterAddress != "" {
		address := strings.Split(clusterAddress, ",")
		baseConfig.Cluster.Servers = address
	}

	if baseConfig.Cluster.IsValid() {
		center.ClusterCenter.ConnectCluster(baseConfig.Cluster)
	} else {
		log.Infof("disable cluster %v", baseConfig.Cluster)
	}

	var moduleNames []string = nil
	moduleConfig := os.Getenv("Module")
	if moduleConfig != "" {
		moduleNames = strings.Split(moduleConfig, ",")
	}

	//logo := `
	//╔═║║  ╝╔═╝╔═ ╔═╝╔═ ╔═║═╔╝
	//╔═║║  ║╔═╝║ ║══║╔═║║ ║ ║
	//╝ ╝══╝╝══╝╝ ╝══╝══ ══╝ ╝
	//`

	f, err := os.Open(configPath + "/logo.txt")
	if err == nil {
		data, _ := ioutil.ReadAll(f)
		fmt.Println(string(data))
	} else {
		log.Debug(err)
	}

	log.Infof("AliensBoot %v starting up...", config.Version)

	//module.Register(database.Module)

	// module
	for i := 0; i < len(mods); i++ {
		if moduleNames != nil && !contains(moduleNames, mods[i].GetName()) {
			log.Warnf("ignore module %v", mods[i].GetName())
			continue
		}
		module.Register(mods[i])
	}

	module.Init()
	// console
	console.Init(baseConfig.ConsolePort, baseConfig.ConsolePrompt)

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Infof("AliensBoot closing down (signal: %v)", sig)
	console.Destroy()
	module.Destroy()
	//close cluster
	center.ClusterCenter.Close()
}

func contains(modules []string, module string) bool {
	for _, moduleName := range modules {
		if strings.TrimSpace(moduleName) == module {
			return true
		}
	}
	return false
}
