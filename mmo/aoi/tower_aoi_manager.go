package aoi

import "github.com/KylinHe/aliensboot-core/mmo/unit"

func NewTowerAOIManager(minX, maxX, minY, maxY unit.Coord, towerRange unit.Coord) Manager {
	this := &TowerAOIManager{minX: minX, maxX: maxX, minY: minY, maxY: maxY, towerRadius: towerRange}
	this.init()
	return this
}

func (this *TowerAOIManager) init() {
	numXSlots := int((this.maxX-this.minX)/this.towerRadius) + 1
	this.xTowerNum = numXSlots
	numYSlots := int((this.maxY-this.minY)/this.towerRadius) + 1
	this.yTowerNum = numYSlots
	this.towers = make([][]tower, numXSlots)
	for i := 0; i < numXSlots; i++ {
		this.towers[i] = make([]tower, numYSlots)
		for j := 0; j < numYSlots; j++ {
			this.towers[i][j].init(i, j)
		}
	}

}

type TowerAOIManager struct {
	minX, maxX, minY, maxY unit.Coord //AOI 监控区域
	towerRadius            unit.Coord //灯塔监控范围
	towers                 [][]tower  //灯塔二位数组
	xTowerNum, yTowerNum   int        //灯塔的横向和纵向的数量
}

//update aoi when change view radius
func (this *TowerAOIManager) ChangeViewRadius(aoi *AOI, newDist unit.Coord) {
	diff := newDist - aoi.viewRadius
	if diff == 0 {
		return
	}
	aoi.viewRadius = newDist
	//TODO 优化计算、通过增量查找新增和删除的灯塔的wathcer
	//视野变大
	if diff > 0 {
		this.visitWatchedTowers(aoi.x, aoi.y, aoi.viewRadius, func(tower *tower) {
			tower.addWatcher(aoi)
		})
	} else {
		this.visitWatchedTowers(aoi.x, aoi.y, aoi.viewRadius, func(tower *tower) {
			tower.removeWatcher(aoi)
		})
	}
}

//add AOI node
func (this *TowerAOIManager) Enter(aoi *AOI, x, y unit.Coord) {
	aoi.x = x
	aoi.y = y

	this.visitWatchedTowers(x, y, aoi.viewRadius, func(tower *tower) {
		tower.addWatcher(aoi)
	})

	t := this.getTower(x, y)
	t.addAOINode(aoi, nil)
}

func (this *TowerAOIManager) Leave(aoi *AOI) {
	aoi.tower.removeAOINode(aoi, true)
	this.visitWatchedTowers(aoi.x, aoi.y, aoi.viewRadius, func(tower *tower) {
		tower.removeWatcher(aoi)
	})
}

func (this *TowerAOIManager) Moved(aoiNode *AOI, x, y unit.Coord) {
	oldX := aoiNode.x
	oldY := aoiNode.y
	aoiNode.x = x
	aoiNode.y = y
	oldTower := aoiNode.tower
	newTower := this.getTower(x, y)

	if oldTower != newTower {
		oldTower.removeAOINode(aoiNode, false)
		newTower.addAOINode(aoiNode, oldTower)
	}

	oxMin, oxMax, oyMin, oyMax := this.getWatchedTowers(oldX, oldY, aoiNode.viewRadius)
	xMin, xMax, yMin, yMax := this.getWatchedTowers(x, y, aoiNode.viewRadius)

	for xi := oxMin; xi <= oxMax; xi++ {
		for yi := oyMin; yi <= oyMax; yi++ {
			if xi >= xMin && xi <= xMax && yi >= yMin && yi <= yMax {
				continue
			}
			tower := &this.towers[xi][yi]
			tower.removeWatcher(aoiNode)
		}
	}

	for xi := xMin; xi <= xMax; xi++ {
		for yi := yMin; yi <= yMax; yi++ {
			if xi >= oxMin && xi <= oxMax && yi >= oyMin && yi <= oyMax {
				continue
			}
			tower := &this.towers[xi][yi]
			tower.addWatcher(aoiNode)
		}
	}
}

func (this *TowerAOIManager) transTowerXY(x, y unit.Coord) (int, int) {
	xi := int((x - this.minX) / this.towerRadius)
	yi := int((y - this.minY) / this.towerRadius)
	return this.normalizeXi(xi), this.normalizeYi(yi)
}

func (this *TowerAOIManager) normalizeXi(xi int) int {
	if xi < 0 {
		xi = 0
	} else if xi >= this.xTowerNum {
		xi = this.xTowerNum - 1
	}
	return xi
}

func (this *TowerAOIManager) normalizeYi(yi int) int {
	if yi < 0 {
		yi = 0
	} else if yi >= this.yTowerNum {
		yi = this.yTowerNum - 1
	}
	return yi
}

func (this *TowerAOIManager) getTower(x, y unit.Coord) *tower {
	xi, yi := this.transTowerXY(x, y)
	return &this.towers[xi][yi]
}

func (this *TowerAOIManager) getWatchedTowers(x unit.Coord, y unit.Coord, aoiDistance unit.Coord) (int, int, int, int) {
	xMin, yMin := this.transTowerXY(x-aoiDistance, y-aoiDistance)
	xMax, yMax := this.transTowerXY(x+aoiDistance, y+aoiDistance)
	return xMin, xMax, yMin, yMax
}

func (this *TowerAOIManager) visitWatchedTowers(x unit.Coord, y unit.Coord, aoiDistance unit.Coord, callback func(*tower)) {
	xMin, xMax, yMin, yMax := this.getWatchedTowers(x, y, aoiDistance)
	for xi := xMin; xi <= xMax; xi++ {
		for yi := yMin; yi <= yMax; yi++ {
			tower := &this.towers[xi][yi]
			callback(tower)
		}
	}
}
