package main

import (
	"net/http"
)

// 오디오 데이터를 스트리밍 받을 사용자가 서버로 요청.
// 클라이언트가 서버로 localhost:8368/stream을 요청.
// 이 요청을 날리면, 클라이언트 <-> 서버가 커넥션을 맺음.
// nginx 설정에 따라 커넥션의 지속시간이 달라짐.
func handleStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "audio/mpeg")
	ch := bc.Register()
	defer bc.Unregister(ch)

	// 채널이 닫힐 때 까지 무한대기
	// 클라이언트가 직접 연결을 종료할 때까지 24시간 커넥션 연결
	for chunk := range ch {
		w.Write(chunk)
		w.(http.Flusher).Flush()
	}
}

// 재생할 오디오를 Bcast 서버에 등록하기 위한 요청
// 사용자 혹은 AI Agent가 localhost:8368/inject에
// raw 오디오 바이너리(mp3 등)를 HTTP Body에 담아 POST 요청
// 흠... 그런데 저 오디오 바이너리를 어떻게, HTTP Body에 담아서 POST 요청을 날리지?
//  1. 정적 오디오 파일 (.mp3)
//  2. 스트리밍 오디오 스트리밍 바이너리 (라디오)
//  3. Google TTS 오디오 스트리밍 바이너리
func handleInject(w http.ResponseWriter, r *http.Request) {
	done := make(chan struct{})

	// r.Body = raw 오디오 바이너리가 흘러들어오는 통로
	// 워커가 이 통로에서 4096바이트씩 읽어서 클라이언트에게 스트리밍
	// done이 닫힐 때까지 HTTP 커넥션 유지 (중간에 끊기면 오디오가 잘림)
	InjectQueue <- InjectItem{reader: r.Body, done: done}

	// 워커가 다 읽을 때까지 대기 = 커넥션 유지
	// 워커에서 close(item.done)으로 채널을 닫으면 블로킹이 풀림
	// 즉 "워커가 다 읽었다고 신호 줄 때까지 HTTP 커넥션 유지해"라는 의미
	<-done
	w.WriteHeader(http.StatusOK)
}
