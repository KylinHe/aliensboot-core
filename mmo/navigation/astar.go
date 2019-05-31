/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/8/27
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package main

import (
	"container/heap"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type OpenList []*_AstarPoint

func (self OpenList) Len() int           { return len(self) }
func (self OpenList) Less(i, j int) bool { return self[i].fVal < self[j].fVal }
func (self OpenList) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }

func (this *OpenList) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*this = append(*this, x.(*_AstarPoint))
}

func (this *OpenList) Pop() interface{} {
	old := *this
	n := len(old)
	x := old[n-1]
	*this = old[0 : n-1]
	return x
}

type _Point struct {
	x    int
	y    int
	view string
}

//========================================================================================

// 保存地图的基本信息
type Map struct {
	points [][]_Point
	blocks map[string]*_Point
	maxX   int
	maxY   int
}

func NewMap(charMap []string) (m Map) {
	m.points = make([][]_Point, len(charMap))
	m.blocks = make(map[string]*_Point, len(charMap)*2)
	for x, row := range charMap {
		cols := strings.Split(row, " ")
		m.points[x] = make([]_Point, len(cols))
		for y, view := range cols {
			m.points[x][y] = _Point{x, y, view}
			if view == "X" {
				m.blocks[pointAsKey(x, y)] = &m.points[x][y]
			}
		} // end of cols
	} // end of row

	m.maxX = len(m.points)
	m.maxY = len(m.points[0])

	return m
}

func (this *Map) getAdjacentPoint(curPoint *_Point) (adjacents []*_Point) {
	if x, y := curPoint.x, curPoint.y-1; x >= 0 && x < this.maxX && y >= 0 && y < this.maxY {
		adjacents = append(adjacents, &this.points[x][y])
	}
	if x, y := curPoint.x+1, curPoint.y-1; x >= 0 && x < this.maxX && y >= 0 && y < this.maxY {
		adjacents = append(adjacents, &this.points[x][y])
	}
	if x, y := curPoint.x+1, curPoint.y; x >= 0 && x < this.maxX && y >= 0 && y < this.maxY {
		adjacents = append(adjacents, &this.points[x][y])
	}
	if x, y := curPoint.x+1, curPoint.y+1; x >= 0 && x < this.maxX && y >= 0 && y < this.maxY {
		adjacents = append(adjacents, &this.points[x][y])
	}
	if x, y := curPoint.x, curPoint.y+1; x >= 0 && x < this.maxX && y >= 0 && y < this.maxY {
		adjacents = append(adjacents, &this.points[x][y])
	}
	if x, y := curPoint.x-1, curPoint.y+1; x >= 0 && x < this.maxX && y >= 0 && y < this.maxY {
		adjacents = append(adjacents, &this.points[x][y])
	}
	if x, y := curPoint.x-1, curPoint.y; x >= 0 && x < this.maxX && y >= 0 && y < this.maxY {
		adjacents = append(adjacents, &this.points[x][y])
	}
	if x, y := curPoint.x-1, curPoint.y-1; x >= 0 && x < this.maxX && y >= 0 && y < this.maxY {
		adjacents = append(adjacents, &this.points[x][y])
	}
	return adjacents
}

func (this *Map) PrintMap(path *SearchRoad) {
	fmt.Println("map's border:", this.maxX, this.maxY)
	for x := 0; x < this.maxX; x++ {
		for y := 0; y < this.maxY; y++ {
			if path != nil {
				if x == path.start.x && y == path.start.y {
					fmt.Print("S")
					goto NEXT
				}
				if x == path.end.x && y == path.end.y {
					fmt.Print("E")
					goto NEXT
				}
				for i := 0; i < len(path.TheRoad); i++ {
					if path.TheRoad[i].x == x && path.TheRoad[i].y == y {
						fmt.Print("*")
						goto NEXT
					}
				}
			}
			fmt.Print(this.points[x][y].view)
		NEXT:
		}
		fmt.Println()
	}
}

func pointAsKey(x, y int) (key string) {
	key = strconv.Itoa(x) + "," + strconv.Itoa(y)
	return key
}

//========================================================================================

type _AstarPoint struct {
	_Point
	father *_AstarPoint
	gVal   int
	hVal   int
	fVal   int
}

func NewAstarPoint(p *_Point, father *_AstarPoint, end *_AstarPoint) (ap *_AstarPoint) {
	ap = &_AstarPoint{*p, father, 0, 0, 0}
	if end != nil {
		ap.calcFVal(end)
	}
	return ap
}

func (this *_AstarPoint) calcGVal() int {
	if this.father != nil {
		deltaX := math.Abs(float64(this.father.x - this.x))
		deltaY := math.Abs(float64(this.father.y - this.y))
		if deltaX == 1 && deltaY == 0 {
			this.gVal = this.father.gVal + 10
		} else if deltaX == 0 && deltaY == 1 {
			this.gVal = this.father.gVal + 10
		} else if deltaX == 1 && deltaY == 1 {
			this.gVal = this.father.gVal + 14
		} else {
			panic("father point is invalid!")
		}
	}
	return this.gVal
}

func (this *_AstarPoint) calcHVal(end *_AstarPoint) int {
	this.hVal = int(math.Abs(float64(end.x-this.x)) + math.Abs(float64(end.y-this.y)))
	return this.hVal
}

func (this *_AstarPoint) calcFVal(end *_AstarPoint) int {
	this.fVal = this.calcGVal() + this.calcHVal(end)
	return this.fVal
}

//========================================================================================

type SearchRoad struct {
	theMap  *Map
	start   _AstarPoint
	end     _AstarPoint
	closeLi map[string]*_AstarPoint
	openLi  OpenList
	openSet map[string]*_AstarPoint
	TheRoad []*_AstarPoint
}

func NewSearchRoad(startx, starty, endx, endy int, m *Map) *SearchRoad {
	sr := &SearchRoad{}
	sr.theMap = m
	sr.start = *NewAstarPoint(&_Point{startx, starty, "S"}, nil, nil)
	sr.end = *NewAstarPoint(&_Point{endx, endy, "E"}, nil, nil)
	sr.TheRoad = make([]*_AstarPoint, 0)
	sr.openSet = make(map[string]*_AstarPoint, m.maxX+m.maxY)
	sr.closeLi = make(map[string]*_AstarPoint, m.maxX+m.maxY)

	heap.Init(&sr.openLi)
	heap.Push(&sr.openLi, &sr.start) // 首先把起点加入开放列表
	sr.openSet[pointAsKey(sr.start.x, sr.start.y)] = &sr.start
	// 将障碍点放入关闭列表
	for k, v := range m.blocks {
		sr.closeLi[k] = NewAstarPoint(v, nil, nil)
	}

	return sr
}

func (this *SearchRoad) FindoutRoad() bool {
	for len(this.openLi) > 0 {
		// 将节点从开放列表移到关闭列表当中。
		x := heap.Pop(&this.openLi)
		curPoint := x.(*_AstarPoint)
		delete(this.openSet, pointAsKey(curPoint.x, curPoint.y))
		this.closeLi[pointAsKey(curPoint.x, curPoint.y)] = curPoint

		//fmt.Println("curPoint :", curPoint.x, curPoint.y)
		adjacs := this.theMap.getAdjacentPoint(&curPoint._Point)
		for _, p := range adjacs {
			//fmt.Println("\t adjact :", p.x, p.y)
			theAP := NewAstarPoint(p, curPoint, &this.end)
			if pointAsKey(theAP.x, theAP.y) == pointAsKey(this.end.x, this.end.y) {
				// 找出路径了, 标记路径
				for theAP.father != nil {
					this.TheRoad = append(this.TheRoad, theAP)
					theAP.view = "*"
					theAP = theAP.father
				}
				return true
			}

			_, ok := this.closeLi[pointAsKey(p.x, p.y)]
			if ok {
				continue
			}

			existAP, ok := this.openSet[pointAsKey(p.x, p.y)]
			if !ok {
				heap.Push(&this.openLi, theAP)
				this.openSet[pointAsKey(theAP.x, theAP.y)] = theAP
			} else {
				oldGVal, oldFather := existAP.gVal, existAP.father
				existAP.father = curPoint
				existAP.calcGVal()
				// 如果新的节点的G值还不如老的节点就恢复老的节点
				if existAP.gVal > oldGVal {
					// restore father
					existAP.father = oldFather
					existAP.gVal = oldGVal
				}
			}

		}
	}

	return false
}

//========================================================================================

func main() {
	presetMap := []string{
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		"X . X X X X X X X X X X X X X X X X X X X X X X X X X",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		"X X X X X X X X X X X X X X X X X X X X X X X X . X X",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
		". . . . . . . . . . . . . . . . . . . . . . . . . . .",
	}
	m := NewMap(presetMap)
	m.PrintMap(nil)

	searchRoad := NewSearchRoad(0, 0, 18, 10, &m)
	if searchRoad.FindoutRoad() {
		fmt.Println("找到了， 你看！")
		m.PrintMap(searchRoad)
	} else {
		fmt.Println("找不到路径！")
	}
}
