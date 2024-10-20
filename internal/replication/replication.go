package replication

import (
	"fmt"
	"math/rand"
	"sync"

	"dfs/config"
	"dfs/internal/peer"
)

type ReplicationManager struct {
	cfg           *config.Config
	chunkLocation map[string][]string // chunkID -> []peerID
	mu            sync.RWMutex
}

func NewReplicationManager(cfg *config.Config) *ReplicationManager {
	return &ReplicationManager{
		cfg:           cfg,
		chunkLocation: make(map[string][]string),
	}
}

func (rm *ReplicationManager) ReplicateChunk(chunk []byte, chunkID string, peers map[string]*peer.Peer) error {
	targetPeers := rm.selectPeers(peers, rm.cfg.ReplicationFactor)

	var wg sync.WaitGroup
	errChan := make(chan error, len(targetPeers))

	for _, p := range targetPeers {
		wg.Add(1)
		go func(p *peer.Peer) {
			defer wg.Done()
			if err := p.SendChunk(chunk, chunkID); err != nil {
				errChan <- fmt.Errorf("failed to send chunk to peer %s: %v", p.ID(), err)
			}
		}(p)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	rm.mu.Lock()
	rm.chunkLocation[chunkID] = append(rm.chunkLocation[chunkID], rm.getPeerIDs(targetPeers)...)
	rm.mu.Unlock()

	return nil
}

func (rm *ReplicationManager) GetChunkLocations(chunkID string) []string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.chunkLocation[chunkID]
}

func (rm *ReplicationManager) selectPeers(peers map[string]*peer.Peer, count int) []*peer.Peer {
	peerList := make([]*peer.Peer, 0, len(peers))
	for _, p := range peers {
		peerList = append(peerList, p)
	}

	rand.Shuffle(len(peerList), func(i, j int) {
		peerList[i], peerList[j] = peerList[j], peerList[i]
	})

	if len(peerList) < count {
		count = len(peerList)
	}

	return peerList[:count]
}

func (rm *ReplicationManager) getPeerIDs(peers []*peer.Peer) []string {
	ids := make([]string, len(peers))
	for i, p := range peers {
		ids[i] = p.ID()
	}
	return ids
}