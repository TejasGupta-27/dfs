package peer

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	conn     net.Conn
	id       string
	sendChan chan []byte
	mu       sync.Mutex
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		conn:     conn,
		id:       conn.RemoteAddr().String(),
		sendChan: make(chan []byte, 100),
	}
}

func (p *Peer) ID() string {
	return p.id
}

func (p *Peer) Handle() {
	go p.readLoop()
	go p.writeLoop()
}

func (p *Peer) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}
	close(p.sendChan)
}

func (p *Peer) SendChunk(chunk []byte, chunkID string) error {
	msg := map[string]interface{}{
		"type":    "chunk",
		"chunkID": chunkID,
		"data":    chunk,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	p.sendChan <- data
	return nil
}

func (p *Peer) readLoop() {
	defer p.Close()
	for {
		var msg map[string]interface{}
		decoder := json.NewDecoder(p.conn)
		if err := decoder.Decode(&msg); err != nil {
			fmt.Printf("Error reading from peer %s: %v\n", p.id, err)
			return
		}
		// Process received message
		fmt.Printf("Received message from peer %s: %v\n", p.id, msg)
	}
}

func (p *Peer) writeLoop() {
	for data := range p.sendChan {
		_, err := p.conn.Write(data)
		if err != nil {
			fmt.Printf("Error writing to peer %s: %v\n", p.id, err)
			return
		}
	}
}