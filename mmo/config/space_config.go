/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/8/25
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package config

import "github.com/KylinHe/aliensboot-core/mmo/unit"

type SpaceConfig struct {
	Type       string     //类型名 用来分类
	MinX       unit.Coord //x轴左边距
	MaxX       unit.Coord
	MinY       unit.Coord
	MaxY       unit.Coord
	TowerRange unit.Coord
}

type EntityConfig struct {
}
