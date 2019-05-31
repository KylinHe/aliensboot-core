/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2018/7/9
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package cipher

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	md5Hash := hex.EncodeToString(h.Sum(nil))
	return md5Hash
}

