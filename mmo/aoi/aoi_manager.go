package aoi

import "github.com/KylinHe/aliensboot-core/mmo/unit"

type Manager interface {
	ChangeViewRadius(aoi *AOI, radius unit.Coord)
	Enter(aoi *AOI, x, y unit.Coord)
	Leave(aoi *AOI)
	Moved(aoi *AOI, x, y unit.Coord)
}
