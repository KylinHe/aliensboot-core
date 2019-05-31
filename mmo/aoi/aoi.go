package aoi

import "github.com/KylinHe/aliensboot-core/mmo/unit"

type xzAOI struct {
	aoi          *AOI
	neighbors    map[*xzAOI]struct{}
	xPrev, xNext *xzAOI
	yPrev, yNext *xzAOI
	markVal      int
}

type AOI struct {
	x          unit.Coord
	y          unit.Coord
	viewRadius unit.Coord //视野范围

	tower *tower

	Callback AOICallback

	implData interface{}
	//implData interface{}
}

func (this *AOI) GetViewRadius() unit.Coord {
	return this.viewRadius
}

func NewAOI(data AOICallback, viewRadius unit.Coord) *AOI {
	return &AOI{
		viewRadius: viewRadius,
		Callback:   data,
	}
}

type AOICallback interface {
	OnEnterAOI(other *AOI)
	OnLeaveAOI(other *AOI)
}
