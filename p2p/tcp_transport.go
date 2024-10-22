package p2p

import (
	"fmt"
	"net"
	"sync"

	
)

// Tcp peer - remote node established over a TCP established connection

type TCPPeer struct{
	conn net.Conn
	// if we dial and retrieve a conn => outbound ==true 
	// if we accept and retrieve a conn => outbound == false 

	outbound bool 

}
func NewTCPPeer(conn net.Conn,outbound bool)*TCPPeer{
	return &TCPPeer{
		conn : conn ,
		outbound: outbound,
	}
}

type TCPTransportOpts struct{
	listenAddr string 
	HandshakeFunc HandshakeFunc
	Decoder	Decoder
	

}


type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	mu sync.RWMutex
	peers map[net.Addr] Peer
}

func NewTCPTransport(opts TCPTransportOpts)*TCPTransport{
	return &TCPTransport{
		TCPTransportOpts: opts,
	}

}

func (t*TCPTransport)ListenAndAccept() error {
	var err error 
	t.listener,err = net.Listen("tcp",t.listenAddr)
	if err!=nil{
		return err
	}
	go t.startAcceptLoop()
	return nil

}
func(t*TCPTransport) startAcceptLoop(){
	for {
		conn,err :=t.listener.Accept()
		if err!=nil{
			fmt.Printf("TCP accept error : %s\n",err)
		}
		
		go t.handleConn(conn)
}
}
type Temp struct{}
func (t*TCPTransport) handleConn(conn net.Conn) {
	peer :=NewTCPPeer(conn,true)
	if err := t.HandshakeFunc(peer);err!=nil{
		conn.Close()
		fmt.Printf("TCP Handshake error : %s\n",err)
		return 

	}
	
	msg := &Temp{}
	for{
		if err:=t.Decoder.Decode(conn,msg);err!=nil{
			fmt.Printf("TCP error : %s\n",err)
			continue
		}

	}				
	
}

