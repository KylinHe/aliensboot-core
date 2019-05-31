package gate

import (
	"github.com/KylinHe/aliensboot-core/chanrpc"
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/KylinHe/aliensboot-core/network"
	"net"
	"reflect"
)

type Gate struct {
	Processor    network.Processor
	AgentChanRPC *chanrpc.Server

	TcpConfig config.TCPConfig
	WsConfig  config.WsConfig
	//UdpConfig config.UDPConfig
	KcpConfig config.KCPConfig
}

func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if gate.WsConfig.Address != "" {
		wsServer = new(network.WSServer)
		wsServer.WsConfig = gate.WsConfig
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go(CommandAgentNew, a)
			}
			return a
		}
	}

	var tcpServer *network.TCPServer
	if gate.TcpConfig.Address != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.TCPConfig = gate.TcpConfig
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go(CommandAgentNew, a)
			}
			return a
		}
	}

	var kcpServer *network.KCPServer
	if gate.KcpConfig.Address != "" {
		kcpServer = new(network.KCPServer)
		kcpServer.KCPConfig = gate.KcpConfig
		kcpServer.NewAgent = func(conn network.Conn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}

	}

	if kcpServer != nil {
		kcpServer.Start()
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if kcpServer != nil {
		kcpServer.Close()
	}

	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

type agent struct {
	conn     network.Conn
	gate     *Gate
	userData interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debugf("read message: %v", err)
			break
		}

		if a.gate.Processor != nil {
			msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				log.Debugf("unmarshal message error: %v", err)
				break
			}
			//msgType := reflect.TypeOf(msg)
			a.gate.AgentChanRPC.Go(CommandAgentMsg, msg, a)
			//err = a.gate.Processor.Route(msg, a)
			//if err != nil {
			//	log.Debug("route message error: %v", err)
			//	break
			//}
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil {
		err := a.gate.AgentChanRPC.Call0(CommandAgentClose, a)
		if err != nil {
			log.Errorf("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteMsg(msg interface{}) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Errorf("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Errorf("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
