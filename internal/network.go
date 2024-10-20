package internal

import (
	"fmt"
	"net"
	"sync"

	"github.com/TejasGupta-27/dfs/config"
	"github.com/TejasGupta-27/dfs/internal/peer"
)

type Network struct {
	cfg      *config.Config
	listener net.Listener
	peers    map[string]*Peer
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
		peers:    make(map[string]*Peer),
	}, nil
}

func (n *Network) Start() {
	go n.acceptConnections()
}

func (n *Network) Stop() {
	n.listener.Close()
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, peer := range n.peers {
		peer.Close()
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
	peer := NewPeer(conn)
	n.mu.Lock()
	n.peers[peer.ID()] = peer
	n.mu.Unlock()

	peer.Handle()

	n.mu.Lock()
	delete(n.peers, peer.ID())
	n.mu.Unlock()
}