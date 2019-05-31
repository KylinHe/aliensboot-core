package network

import (
	"fmt"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/xtaci/kcp-go"
	"net"
	"runtime"
	"sync"
)

type KCPConn struct {
	sync.Mutex
	conn      *kcp.UDPSession
	writeChan chan []byte
	closeFlag bool
	msgParser *MsgParser
}

func newKCPConn(conn *kcp.UDPSession, pendingWriteNum int, msgParser *MsgParser) *KCPConn {
	kcpConn := new(KCPConn)
	kcpConn.conn = conn
	kcpConn.writeChan = make(chan []byte, pendingWriteNum)
	kcpConn.msgParser = msgParser

	go func() {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				stackInfo := fmt.Sprintf("%s", buf[:n])
				log.Error("panic stack info %s", stackInfo)
			}
		}()
		for b := range kcpConn.writeChan {
			if b == nil {
				break
			}

			_, err := conn.Write(b)
			if err != nil {
				break
			}
		}

		conn.Close()
		kcpConn.Lock()
		kcpConn.closeFlag = true
		kcpConn.Unlock()
	}()

	return kcpConn
}

func (kcpConn *KCPConn) doDestroy() {
	//kcpConn.conn.SetLinger(0)
	kcpConn.conn.Close()

	if !kcpConn.closeFlag {
		close(kcpConn.writeChan)
		kcpConn.closeFlag = true
	}
}

func (kcpConn *KCPConn) Destroy() {
	kcpConn.Lock()
	defer kcpConn.Unlock()

	kcpConn.doDestroy()
}

func (kcpConn *KCPConn) Close() {
	kcpConn.Lock()
	defer kcpConn.Unlock()
	if kcpConn.closeFlag {
		return
	}

	kcpConn.doWrite(nil)
	kcpConn.closeFlag = true
}

func (kcpConn *KCPConn) doWrite(b []byte) {
	if len(kcpConn.writeChan) == cap(kcpConn.writeChan) {
		log.Debug("close conn: channel full")
		kcpConn.doDestroy()
		return
	}

	kcpConn.writeChan <- b
}

// b must not be modified by the others goroutines
func (kcpConn *KCPConn) Write(b []byte) (int, error) {
	kcpConn.Lock()
	defer kcpConn.Unlock()
	if kcpConn.closeFlag || b == nil {
		return -1, nil
	}

	kcpConn.doWrite(b)
	return -1, nil
}

func (kcpConn *KCPConn) Read(b []byte) (int, error) {
	return kcpConn.conn.Read(b)
}

func (kcpConn *KCPConn) LocalAddr() net.Addr {
	return kcpConn.conn.LocalAddr()
}

func (kcpConn *KCPConn) RemoteAddr() net.Addr {
	return kcpConn.conn.RemoteAddr()
}


func (kcpConn *KCPConn) ReadMsg() ([]byte, error) {
	return kcpConn.msgParser.Read(kcpConn)
}

func (kcpConn *KCPConn) WriteMsg(args ...[]byte) error {
	return kcpConn.msgParser.Write(kcpConn, args...)
}

//
//func (kcpConn *KCPConn) ReadMsg() ([]byte, error) {
//	buf := make([]byte, 65536)
//
//	n, err := kcpConn.conn.Read(buf)
//	return buf[2:n], err
//}
//
//func (kcpConn *KCPConn) WriteMsg(args ...[]byte) error {
//	if len(args) == 1 {
//		_, err := kcpConn.conn.Write(args[0])
//		return err
//	}
//	return errors.New("msg packet length must be 1")
//}