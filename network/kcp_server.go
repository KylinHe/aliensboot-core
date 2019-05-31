package network

import (
	"crypto/sha1"
	"fmt"
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"net"
	"runtime"
	"sync"
	"time"
)

type KCPServer struct {
	config.KCPConfig

	NewAgent   func(conn Conn) Agent

	conns      ConnSet
	mutexConns sync.Mutex
	ln         *kcp.Listener
	wgLn       sync.WaitGroup
	wgConns    sync.WaitGroup

	msgParser *MsgParser
}

func (server *KCPServer) Start() {
	server.init()
	go server.run()
}

func (server *KCPServer) defaultConfig() {
	if server.DataShard == 0 {
		server.DataShard = 10
	}
	if server.ParityShard == 0 {
		server.ParityShard = 3
	}
	if server.Mode == "" {
		server.Mode = "fast3"
	}
	if server.MTU <= 0 {
		server.MTU = 1350
	}
	if server.SndWnd == 0 {
		server.SndWnd = 1024
	}
	if server.RcvWnd == 0 {
		server.RcvWnd = 1024
	}
	if server.SockBuf == 0 {
		server.SockBuf = 4194304
	}

	switch server.Mode {
		case "normal":
			server.NoDelay, server.Interval, server.Resend, server.NoCongestion = 0, 40, 2, 1
		case "fast":
			server.NoDelay, server.Interval, server.Resend, server.NoCongestion = 0, 30, 2, 1
		case "fast2":
			server.NoDelay, server.Interval, server.Resend, server.NoCongestion = 1, 20, 2, 1
		case "fast3":
			server.NoDelay, server.Interval, server.Resend, server.NoCongestion = 1, 10, 2, 1
	}


}

func (server *KCPServer) init() {
	server.defaultConfig()
	ln, err := kcp.ListenWithOptions(server.Address, server.getBlockCrypt(), server.DataShard, server.ParityShard)

	if err != nil {
		log.Fatalf("%v", err)
	}

	if server.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}

	if err := ln.SetReadBuffer(server.SockBuf); err != nil {
		log.Fatalf("SetReadBuffer:", err)
	}
	if err := ln.SetWriteBuffer(server.SockBuf); err != nil {
		log.Fatalf("SetWriteBuffer:", err)
	}

	server.ln = ln
	server.conns = make(ConnSet)

	// msg parser
	msgParser := NewMsgParser()
	msgParser.SetMsgLen(server.LenMsgLen, server.MinMsgLen, server.MaxMsgLen)
	msgParser.SetByteOrder(server.LittleEndian)
	server.msgParser = msgParser

}

func (server *KCPServer) run() {
	server.wgLn.Add(1)
	defer server.wgLn.Done()

	var tempDelay time.Duration
	for {
		conn, err := server.ln.AcceptKCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Errorf("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		conn.SetStreamMode(true)
		conn.SetMtu(server.MTU)
		conn.SetWindowSize(server.SndWnd, server.RcvWnd)
		conn.SetACKNoDelay(server.AckNodelay)
		conn.SetWriteDelay(false)
		conn.SetNoDelay(server.NoDelay, server.Interval, server.Resend, server.NoCongestion)

		server.mutexConns.Lock()
		if len(server.conns) >= server.MaxConnNum {
			server.mutexConns.Unlock()
			conn.Close()
			log.Debug("too many connections")
			continue
		}
		server.conns[conn] = struct{}{}
		server.mutexConns.Unlock()

		server.wgConns.Add(1)

		kcpConn := newKCPConn(conn, server.PendingWriteNum, server.msgParser)
		agent := server.NewAgent(kcpConn)

		go func() {
			defer func() {
				if err := recover(); err != nil {
					buf := make([]byte, 2048)
					n := runtime.Stack(buf, false)
					stackInfo := fmt.Sprintf("%s", buf[:n])
					log.Errorf("panic stack info %s", stackInfo)
				}
			}()
			agent.Run()
			// cleanup
			kcpConn.Close()
			server.mutexConns.Lock()
			delete(server.conns, conn)
			server.mutexConns.Unlock()
			agent.OnClose()

			server.wgConns.Done()
		}()

	}
}


func (server *KCPServer) getBlockCrypt() kcp.BlockCrypt {
	if server.Crypt == "" {
		return nil
	}
	pass := pbkdf2.Key([]byte(server.Key), []byte(server.Salt), 4096, 32, sha1.New)
	var block kcp.BlockCrypt
	switch server.Crypt {
	case "sm4":
		block, _ = kcp.NewSM4BlockCrypt(pass[:16])
	case "tea":
		block, _ = kcp.NewTEABlockCrypt(pass[:16])
	case "xor":
		block, _ = kcp.NewSimpleXORBlockCrypt(pass)
	case "none":
		block, _ = kcp.NewNoneBlockCrypt(pass)
	case "aes-128":
		block, _ = kcp.NewAESBlockCrypt(pass[:16])
	case "aes-192":
		block, _ = kcp.NewAESBlockCrypt(pass[:24])
	case "blowfish":
		block, _ = kcp.NewBlowfishBlockCrypt(pass)
	case "twofish":
		block, _ = kcp.NewTwofishBlockCrypt(pass)
	case "cast5":
		block, _ = kcp.NewCast5BlockCrypt(pass[:16])
	case "3des":
		block, _ = kcp.NewTripleDESBlockCrypt(pass[:24])
	case "xtea":
		block, _ = kcp.NewXTEABlockCrypt(pass[:16])
	case "salsa20":
		block, _ = kcp.NewSalsa20BlockCrypt(pass)
	default:
		server.Crypt = "aes"
		block, _ = kcp.NewAESBlockCrypt(pass)
	}
	return block
}

func (server *KCPServer) Close() {
	server.ln.Close()
	server.wgLn.Wait()

	server.mutexConns.Lock()
	for conn := range server.conns {
		conn.Close()
	}
	server.conns = nil
	server.mutexConns.Unlock()
	server.wgConns.Wait()
}