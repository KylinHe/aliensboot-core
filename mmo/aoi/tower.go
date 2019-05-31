package aoi

import "github.com/KylinHe/aliensboot-core/log"

type tower struct {
	indexX   int
	indexY   int
	aoiNodes map[*AOI]struct{} //当前灯塔范围内的aoi对象
	watchers map[*AOI]struct{} //灯塔在视野范围内的所有aoi
}

func (t *tower) init(indexX int, indexY int) {
	t.indexX = indexX
	t.indexY = indexY
	t.aoiNodes = map[*AOI]struct{}{}
	t.watchers = map[*AOI]struct{}{}
	log.Infof("tower init %v-%v", indexX, indexY)
}

func (t *tower) addAOINode(aoiNode *AOI, fromOtherTower *tower) {
	aoiNode.tower = t
	//log.Info("进入 %v-%v 观测台", t.indexX, t.indexY)
	t.aoiNodes[aoiNode] = struct{}{}
	if fromOtherTower == nil {
		for watcher := range t.watchers {
			if watcher == aoiNode {
				continue
			}
			watcher.Callback.OnEnterAOI(aoiNode)
		}
	} else {
		// aoiNode moved from other tower to this tower
		for watcher := range fromOtherTower.watchers {
			if watcher == aoiNode {
				continue
			}
			if _, ok := t.watchers[watcher]; ok {
				continue
			}
			watcher.Callback.OnLeaveAOI(aoiNode)
		}
		for watcher := range t.watchers {
			if watcher == aoiNode {
				continue
			}
			if _, ok := fromOtherTower.watchers[watcher]; ok {
				continue
			}
			watcher.Callback.OnEnterAOI(aoiNode)
		}
	}
}

func (t *tower) removeAOINode(aoiNode *AOI, notifyWatchers bool) {
	aoiNode.tower = nil
	//log.Info("退出 %v-%v 观测台", t.indexX, t.indexY)
	delete(t.aoiNodes, aoiNode)
	if notifyWatchers {
		for watcher := range t.watchers {
			if watcher == aoiNode {
				continue
			}
			watcher.Callback.OnLeaveAOI(aoiNode)
		}
	}
}

func (t *tower) addWatcher(aoiNode *AOI) {
	if _, ok := t.watchers[aoiNode]; ok {
		log.Info("duplicate add watcher")
		return
	}
	//log.Info("观测台 %v-%v 进入可视范围", t.indexX, t.indexY)
	t.watchers[aoiNode] = struct{}{}
	// now aoiNode can see all aoiNodes under this tower
	for neighbor := range t.aoiNodes {
		if neighbor == aoiNode {
			continue
		}
		aoiNode.Callback.OnEnterAOI(neighbor)
	}
}

func (t *tower) removeWatcher(aoiNode *AOI) {
	if _, ok := t.watchers[aoiNode]; !ok {
		log.Info("duplicate remove watcher")
		return
	}
	//log.Info("观测台 %v-%v 退出可视范围", t.indexX, t.indexY)
	delete(t.watchers, aoiNode)
	for neighbor := range t.aoiNodes {
		if neighbor == aoiNode {
			continue
		}
		aoiNode.Callback.OnLeaveAOI(neighbor)
	}
}
