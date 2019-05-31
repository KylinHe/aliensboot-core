/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/5/10
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package conf

const (
	AnalysisFlag = false               //是否开启性能分析
	Game         = "aliens"            //统计日志索引前缀
	Local        = true                //是否存储到本地
	LocalPrefix  = "aliens_statistics" //本地日志存储的前缀
)

var Config struct {
	ES ESConfig
}

type ESConfig struct {
	Url      string
	Host     string
	Username string
	Password string
}
