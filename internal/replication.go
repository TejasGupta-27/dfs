package internal

import (
	"fmt"
	"math/rand"
	"sync"
	"github.com/TejasGupta-27/dfs/config"
	
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

func (rm *ReplicationManager) ReplicateChunk(chunk Chunk, peers map[string]*Peer) error {
	targetPeers := rm.selectPeers(peers, rm.cfg.ReplicationFactor)

	var wg sync.WaitGroup
	errChan := make(chan error, len(targetPeers))

	for _, peer := range targetPeers {
		wg.Add(1)
		go func(p *Peer) {
			defer wg.Done()
			if err := p.SendChunk(chunk); err != nil {
				errChan <- fmt.Errorf("failed to send chunk to peer %s: %v", p.ID(), err)
			}
		}(peer)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	rm.mu.Lock()
	rm.chunkLocation[chunk.ID] = append(rm.chunkLocation[chunk.ID], rm.getPeerIDs(targetPeers)...)
	rm.mu.Unlock()

	return nil
}

func (rm *ReplicationManager) GetChunkLocations(chunkID string) []string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.chunkLocation[chunkID]
}

func (rm *ReplicationManager) selectPeers(peers map[string]*Peer, count int) []*Peer {
	peerList := make([]*Peer, 0, len(peers))
	for _, peer := range peers {
		peerList = append(peerList, peer)
	}

	rand.Shuffle(len(peerList), func(i, j int) {
		peerList[i], peerList[j] = peerList[j], peerList[i]
	})

	if len(peerList) < count {
		count = len(peerList)
	}

	return peerList[:count]
}

func (rm *ReplicationManager) getPeerIDs(peers []*Peer) []string {
	ids := make([]string, len(peers))
	for i, peer := range peers {
		ids[i] = peer.ID()
	}
	return ids
}