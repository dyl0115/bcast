package main

import "sync"

type Broadcaster struct {
	mu      sync.Mutex
	clients map[chan []byte]struct{}
}

func (b *Broadcaster) Register() chan []byte {
	ch := make(chan []byte, 256)
	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *Broadcaster) Unregister(ch chan []byte) {
	b.mu.Lock()
	delete(b.clients, ch)
	close(ch)
	b.mu.Unlock()
}

func (b *Broadcaster) Send(chunk []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.clients {
		buf := make([]byte, len(chunk))
		copy(buf, chunk)
		ch <- buf
	}
}

var bc = &Broadcaster{clients: make(map[chan []byte]struct{})}
