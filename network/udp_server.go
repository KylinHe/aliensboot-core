/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/6/8
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package network

import (
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/log"
	"net"
)

type UDPServer struct {
	conn *net.UDPConn

	config.UDPConfig

	agents map[string]*UDPAgent

	handle func(data []byte, addr *UDPAgent)
}

func (this *UDPServer) Start() {
	this.init()
	go this.run()
}

func (this *UDPServer) init() {
	this.agents = make(map[string]*UDPAgent)

	udpAddr, err := net.ResolveUDPAddr(this.Protocol, this.Address)
	if err != nil {
		log.Fatalf("star udp server error : %v", err)
	}
	udpConn, err2 := net.ListenUDP(this.Protocol, udpAddr)
	if err2 != nil {
		log.Fatalf("star udp server error : %v", err2)
	}
	this.conn = udpConn
}

func (this *UDPServer) run() {
	defer this.conn.Close()
	for {
		buf := make([]byte, this.MaxMsgLen)
		//读取数据
		len, udpAddr, err := this.conn.ReadFromUDP(buf)
		if err != nil {
			log.Errorf("read udp msg err : %v", err)
		}

		agent := this.agents[udpAddr.String()]
		if agent == nil {
			agent = &UDPAgent{conn: this.conn, udpAddr: udpAddr}
			this.agents[udpAddr.String()] = agent
		}
		this.handle(buf[0:len], agent)
	}
}

func (this *UDPServer) Close() {
	if this.conn != nil {
		this.conn.Close()
	}
}
