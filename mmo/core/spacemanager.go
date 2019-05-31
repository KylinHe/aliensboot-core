/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/3/23
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

var (
	SpaceManager = newSpaceManager()
)

type _SpaceManager struct {
	spaces map[EntityID]*Space
}

func newSpaceManager() *_SpaceManager {
	return &_SpaceManager{
		spaces: map[EntityID]*Space{},
	}
}

func (spmgr *_SpaceManager) GetSpace(id EntityID) *Space {
	return spmgr.spaces[id]
}

func (spmgr *_SpaceManager) putSpace(space *Space) {
	spmgr.spaces[space.GetID()] = space
}

func (spmgr *_SpaceManager) delSpace(id EntityID) {
	delete(spmgr.spaces, id)
}

//release exist space
func (spmgr *_SpaceManager) releaseSpace(id EntityID) {
	delete(spmgr.spaces, id)
}
