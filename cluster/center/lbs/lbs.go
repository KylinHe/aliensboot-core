/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/11/3
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 * Desc:
 *     Load Balance Strategy
 *******************************************************************************/
package lbs

const (
	StrategyPolling  string = "polling"  //轮询
	StrategyHashring string = "hashring" //一致性hash
)

func GetLBS(lbs string) Strategy {
	if lbs == StrategyPolling {
		return NewPollingLBS()
	} else if lbs == StrategyHashring {
		return NewHashRing(400)
	} else {
		return NewPollingLBS()
	}
}

type Strategy interface {

	//Init(nodes []string) //更新所有的负载节点信息

	AddNode(nodeKey string, weight int)

	RemoveNode(nodeKey string)

	GetNode(key string) string //分配可用节点

}
