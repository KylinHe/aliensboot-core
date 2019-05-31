/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/3/21
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package unit

import (
	"fmt"
	"math"
)

var EmptyVector = Vector{0, 0, 0}

// Yaw is the type of entity Yaw
type Yaw float32

// Coord is the of coordinations entity position (x, y, z)
type Coord float32

// Vector is type of entity position
type Vector struct {
	X Coord
	Y Coord
	Z Coord
}

func (p Vector) String() string {
	return fmt.Sprintf("(%.2f, %.2f, %.2f)", p.X, p.Y, p.Z)
}

// DistanceTo calculates distance between two positions
func (p Vector) DistanceTo(o Vector) Coord {
	dx := p.X - o.X
	dy := p.Y - o.Y
	dz := p.Z - o.Z
	return Coord(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// Sub calculates Vector p - Vector o
func (p Vector) Sub(o Vector) Vector {
	return Vector{p.X - o.X, p.Y - o.Y, p.Z - o.Z}
}

func (p Vector) Add(o Vector) Vector {
	return Vector{p.X + o.X, p.Y + o.Y, p.Z + o.Z}
}

// Mul calculates Vector p * m
func (p Vector) Mul(m Coord) Vector {
	return Vector{p.X * m, p.Y * m, p.Z * m}
}

// DirToYaw convert direction represented by Vector to Yaw
func (dir Vector) DirToYaw() Yaw {
	dir.Normalize()

	yaw := math.Acos(float64(dir.X))
	if dir.Z < 0 {
		yaw = math.Pi*2 - yaw
	}

	yaw = yaw / math.Pi * 180 // convert to angle

	if yaw <= 90 {
		yaw = 90 - yaw
	} else {
		yaw = 90 + (360 - yaw)
	}

	return Yaw(yaw)
}

func (p *Vector) Normalize() {
	d := Coord(math.Sqrt(float64(p.X*p.X + p.Y + p.Y + p.Z*p.Z)))
	if d == 0 {
		return
	}
	p.X /= d
	p.Y /= d
	p.Z /= d
}

func (p Vector) Normalized() Vector {
	p.Normalize()
	return p
}
