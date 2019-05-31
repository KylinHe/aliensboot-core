package aoi

import "github.com/KylinHe/aliensboot-core/mmo/unit"

// XZListAOIManager is an implementation of AOICalculator using XZ lists
type XZListAOIManager struct {
	//aoidist    float32
	xSweepList *xAOIList
	zSweepList *yAOIList
}

// NewXZListAOIManager creates a new XZListAOIManager
func NewXZListAOIManager(aoidist unit.Coord) Manager {
	return &XZListAOIManager{
		//aoidist:    aoidist,
		xSweepList: newXAOIList(aoidist),
		zSweepList: newYAOIList(aoidist),
	}
}

// Enter is called when Entity enters Space
func (aoiman *XZListAOIManager) Enter(aoi *AOI, x, y unit.Coord) {
	//aoi.viewRadius = aoiman.aoidist

	xzaoi := &xzAOI{
		aoi:       aoi,
		neighbors: map[*xzAOI]struct{}{},
	}
	aoi.x, aoi.y = x, y
	aoi.implData = xzaoi
	aoiman.xSweepList.Insert(xzaoi)
	aoiman.zSweepList.Insert(xzaoi)
	aoiman.adjust(xzaoi)
}

func (aoiman *XZListAOIManager) ChangeViewRadius(aoi *AOI, radius unit.Coord) {
	//TODO 后续支持
}

// Leave is called when Entity leaves Space
func (aoiman *XZListAOIManager) Leave(aoi *AOI) {
	xzaoi := aoi.implData.(*xzAOI)
	aoiman.xSweepList.Remove(xzaoi)
	aoiman.zSweepList.Remove(xzaoi)
	aoiman.adjust(xzaoi)
}

// Moved is called when Entity moves in Space
func (aoiman *XZListAOIManager) Moved(aoi *AOI, x, y unit.Coord) {
	oldX := aoi.x
	oldY := aoi.y
	aoi.x, aoi.y = x, y
	xzaoi := aoi.implData.(*xzAOI)
	if oldX != x {
		aoiman.xSweepList.Move(xzaoi, oldX)
	}
	if oldY != y {
		aoiman.zSweepList.Move(xzaoi, oldY)
	}
	aoiman.adjust(xzaoi)
}

// adjust is called by Entity to adjust neighbors
func (aoiman *XZListAOIManager) adjust(aoi *xzAOI) {
	aoiman.xSweepList.Mark(aoi)
	aoiman.zSweepList.Mark(aoi)
	// AOI marked twice are neighbors
	for neighbor := range aoi.neighbors {
		if neighbor.markVal == 2 {
			// neighbors kept
			neighbor.markVal = -2 // mark this as neighbor
		} else { // markVal < 2
			// was neighbor, but not any more
			delete(aoi.neighbors, neighbor)
			aoi.aoi.Callback.OnLeaveAOI(neighbor.aoi)
			delete(neighbor.neighbors, aoi)
			neighbor.aoi.Callback.OnLeaveAOI(aoi.aoi)
		}
	}

	// travel in X list again to find all new neighbors, whose markVal == 2
	aoiman.xSweepList.GetClearMarkedNeighbors(aoi)
	// travel in Z list again to unmark all
	aoiman.zSweepList.ClearMark(aoi)
}
