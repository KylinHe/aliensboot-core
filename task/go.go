/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2019/9/11
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package task

import "github.com/KylinHe/aliensboot-core/exception"

// 安全运行携程
func SafeGo(task func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				exception.PrintStackDetail(err)
			}
		}()
		task()
	}()
}
