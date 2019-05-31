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

// ListenWithOptions listens for incoming KCP packets addressed to the local address laddr on the network "udp" with packet encryption,
// dataShards, parityShards defines Reed-Solomon Erasure Coding parameters


// Config for server
type KCPConfig struct {
	Address      string `yaml:"address"`
	MaxConnNum   int 	`yaml:"maxconnnum"`
	PendingWriteNum int `yaml:"pendingwritenum"`

	LenMsgLen       int       `yaml:"lenmsglen"`
	MinMsgLen       uint32    `yaml:"minmsglen"`
	MaxMsgLen       uint32    `yaml:"maxmsglen"`
	LittleEndian    bool      `yaml:"littleendian"`

	Target       string `yaml:"target"`
	Key          string `yaml:"key"`
	Salt		 string `yaml:"salt"`
	Crypt        string `yaml:"crypt"`		  //aes, aes-128, aes-192, salsa20, blowfish, twofish, cast5, 3des, tea, xtea, xor, sm4, none
	Mode         string `yaml:"mode"`		  //profiles: fast3, fast2, fast, normal, manual
	MTU          int    `yaml:"mtu"`		  //et maximum transmission unit for UDP packets
	SndWnd       int    `yaml:"sndwnd"`       //set send window size(num of packets)
	RcvWnd       int    `yaml:"rcvwnd"`		  //set receive window size(num of packets)
	DataShard    int    `yaml:"datashard"`    //set reed-solomon erasure coding - datashard
	ParityShard  int    `yaml:"parityshard"`  //set reed-solomon erasure coding - parityshard
	DSCP         int    `yaml:"dscp"`
	NoComp       bool   `yaml:"nocomp"`
	AckNodelay   bool   `yaml:"acknodelay"`
	NoDelay      int    `yaml:"nodelay"`
	Interval     int    `yaml:"interval"`
	Resend       int    `yaml:"resend"`
	NoCongestion int    `yaml:"nc"`
	SockBuf      int    `yaml:"sockbuf"`
	KeepAlive    int    `yaml:"keepalive"`
	SnmpPeriod   int    `yaml:"snmpperiod"`
	Quiet        bool   `yaml:"quiet"`
}

