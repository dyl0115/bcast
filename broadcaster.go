package main

import "sync"

// Bcast 서버에 연결된 클라이언트들을 등록/제거 하는 객체
type Broadcaster struct {
	// Broadcaster 객체를 여러 고루틴이 사용할 때, mu를 정의해야 LOCK을 걸 수 있음
	// 기본적으로 Go에는 자바의 CuncurrentHashMap 같은 동시성을 자동으로 관리해주는 자료구조가 없음.
	// 동시성을 관리하지 않으면 A클라이언트가 clients map을 수정하는 도중, 다른 B클라이언트가 조작하면 엉망이 됨.
	// Go에서는 sync.Mutex로 동시성을 관리하는게 정석.
	mu sync.Mutex

	// Bcast에 연결된 여러 클라이언트들(데스크탑, 스마트폰 등)을 Map에 담아 관리함.
	// InjectQueue에 오디오 데이터가 들어오면,모든 클라이언트들에게 Broadcasting을 함.
	clients map[chan []byte]struct{}
}

// Broadcaster에 클라이언트를 등록함.
// 오직 Braodcaster 객체만 이 함수를 사용할 수 있음.
// byte[] 타입의 채널을 리턴함.
// 이 채널을 통해 서버는 클라이언트에게 오디오 청크 데이터를 스트리밍.
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

		// 클라이언트 채널이 꽉 차면 여기서 블로킹 됨.. 수정이 필요함.
		ch <- buf
	}
}

var bc = &Broadcaster{clients: make(map[chan []byte]struct{})}
