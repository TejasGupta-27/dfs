package network

import (
	"fmt"
	"net"
	"sync"

	"dfs/config"
	"dfs/internal/peer"
)

type Network struct {
	cfg      *config.Config
	listener net.Listener
	peers    map[string]*peer.Peer
	mu       sync.RWMutex
}

func NewNetwork(cfg *config.Config) (*Network, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}

	return &Network{
		cfg:      cfg,
		listener: listener,
		peers:    make(map[string]*peer.Peer),
	}, nil
}

func (n *Network) Start() error {
	go n.acceptConnections()
	return nil
}

func (n *Network) Stop() {
	n.listener.Close()
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, p := range n.peers {
		p.Close()
	}
}

func (n *Network) acceptConnections() {
	for {
		conn, err := n.listener.Accept()
		if err != nil {
			return
		}
		go n.handleConnection(conn)
	}
}

func (n *Network) handleConnection(conn net.Conn) {
	p := peer.NewPeer(conn)
	n.mu.Lock()
	n.peers[p.ID()] = p
	n.mu.Unlock()

	p.Handle()

	n.mu.Lock()
	delete(n.peers, p.ID())
	n.mu.Unlock()
}