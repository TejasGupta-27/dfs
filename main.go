package main

import (
	"fmt"
	"log"

	"github.com/TejasGupta-27/dfs/p2p"
)

func main(){
	tcpOpts := p2p.TCPTransportOpts{
		listenAddr:=":3000",
		HandshakeFunc:=p2p.NOPHandshakeFunc(),
		
	}
	tr:= p2p.NewTCPTransport(tcpOpts)
	if err :=tr.ListenAndAccept();err!=nil{
		log.Fatal(err)
	}
	select{}

}