/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package util

import "time"

const (
	DurationDay time.Duration = 24 * time.Hour
)

//获取当天开始时间
func GetTodayBegin() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetTodayHourTime(hour int) time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, t.Location())
}

func GetEmptyTime() time.Time {
	return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
}

//获取丹田的结束时间
func GetTodayEnd() time.Time {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return tm1.Add(DurationDay)
}

//获取当前的日志字符串打印
func GetCurrentDayStr() string {
	t := time.Now()
	return t.Format("2006-01-02")
}

func GetTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}


func GetSecondDuration(second int) time.Duration {
	return time.Duration(second) * time.Second
}

