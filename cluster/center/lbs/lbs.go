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
	StrategyPolling  string = "Polling"  //轮询
	StrategyHashRing string = "HashRing" //一致性hash
)

func GetLBS(lbs string) Strategy {
	if lbs == StrategyPolling {
		return NewPollingLBS()
	} else if lbs == StrategyHashRing {
		return NewHashRing(DefaultVirtualSpots)
	} else {
		return NewPollingLBS()
	}
}

func ValidateLBS(lbs string) bool {
	return lbs == StrategyPolling || lbs == StrategyHashRing || lbs == ""
}

type Strategy interface {

	//Init(nodes []string) //更新所有的负载节点信息

	AddNode(nodeKey string, weight int)

	RemoveNode(nodeKey string)

	GetNode(key string) string //分配可用节点

}
