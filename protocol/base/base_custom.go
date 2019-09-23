/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2019/9/23
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package base

func (m *Any) AddHeader(key string, data []byte) {
	if m == nil {
		return
	}
	if m.Header == nil {
		m.Header = make(map[string][]byte, 1)
	}
	m.Header[key] = data
}

func (m *Any) GetHeaderByKey(key string) []byte {
	if m == nil || m.Header == nil {
		return nil
	}
	return m.Header[key]
}

func (m *Any) GetHeaderStrByKey(key string) string {
	if m == nil || m.Header == nil {
		return ""
	}
	data := m.Header[key]
	if data == nil {
		return ""
	}
	return string(data)
}
