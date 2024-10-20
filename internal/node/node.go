package node

import (
	"dfs/config"
	"dfs/internal/network"
	"dfs/internal/file"
	"dfs/internal/replication"
	"dfs/internal/security"
)

type Node struct {
	config      *config.Config
	network     *network.Network
	fileSystem  *file.FileSystem
	replication *replication.ReplicationManager
	security    *security.Security
}

func NewNode(cfg *config.Config) (*Node, error) {
	net, err := network.NewNetwork(cfg)
	if err != nil {
		return nil, err
	}

	sec, err := security.NewSecurity([]byte(cfg.EncryptionKey))
	if err != nil {
		return nil, err
	}

	return &Node{
		config:      cfg,
		network:     net,
		fileSystem:  file.NewFileSystem(cfg),
		replication: replication.NewReplicationManager(cfg),
		security:    sec,
	}, nil
}

func (n *Node) Start() error {
	return n.network.Start()
}

func (n *Node) Stop() {
	n.network.Stop()
}