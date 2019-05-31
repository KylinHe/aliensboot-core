/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/10/25
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package config

type TCPConfig struct {
	Address         string
	MaxConnNum      int
	PendingWriteNum int
	LenMsgLen       int
	MinMsgLen       uint32
	MaxMsgLen       uint32
	LittleEndian    bool
}
