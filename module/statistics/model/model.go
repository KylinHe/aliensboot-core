/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/7/30
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package model

type CallInfo struct {
	count    int32   //调用次数
	interval float64 //调用时间总长
}

func (this *CallInfo) AddCall(interval float64) {
	this.count++
	this.interval += interval
}

func (this *CallInfo) IsEmpty() bool {
	return this.count == 0 || this.interval == 0
}

func (this *CallInfo) DumpData() (bool, int32, float64) {
	if this.IsEmpty() {
		return false, 0, 0
	}
	avg := this.interval / float64(this.count)
	count := this.count

	this.count = 0
	this.interval = 0

	return true, count, avg
}
