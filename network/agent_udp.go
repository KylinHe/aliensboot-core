/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/6/11
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package network

import "net"

type UDPAgent struct {
	conn *net.UDPConn

	udpAddr *net.UDPAddr

	userData interface{}
}

//发送数据
func (this *UDPAgent) WriteMsg(msg interface{}) {
	data, ok := msg.([]byte)
	if ok {
		this.conn.WriteToUDP(data, this.udpAddr)
	}
}

func (this *UDPAgent) SetUserData(userData interface{}) {
	this.userData = userData
}

func (this *UDPAgent) UserData() interface{} {
	return this.userData
}

func (this *UDPAgent) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *UDPAgent) RemoteAddr() net.Addr {
	return this.udpAddr
}

func (this *UDPAgent) Close(userData interface{}) {
	this.userData = userData
}

func (this *UDPAgent) Destroy(userData interface{}) {
	this.userData = userData
}

func (this *UDPAgent) GetID() string {
	if this.udpAddr == nil {
		return ""
	}
	return this.udpAddr.String()
}
