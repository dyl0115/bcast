package main

import "io"

// 여러개의 오디오가 동시에 Queue로 들어오면, 사용자가 뒤죽박죽 섞인 오디오를 듣게 됨.
// 때문에 요청된 여러 오디오(엄밀하게는 데이터 통로)를 순서대로 InjectQueue(고채널)에 넣어줌.
// 백그라운드에서 항상 돌아가는 고루틴이 InjectQueue에서 오디오 청크를 꺼내서 BroadCast.send()로 클라이언트에게 오디오를 스트리밍.
// 최대 20개의 데이터 통로까지 인큐 가능

// InjectItem: inject 요청 하나를 표현하는 구조체
// reader: ffmpeg stdout 또는 파일 등 오디오 소스 (r.Body 직접 사용)
// done: 워커가 읽기를 완료했을 때 닫히는 채널. 핸들러가 done 신호가 올 때까지 HTTP 커넥션을 유지함
type InjectItem struct {

	// io.Reader는 인터페이스.
	// f, _ : = os.Open("song.mp3")에서 f 자체는 데이터가 아님. f는 song.mp3 파일을 읽을 수 있는 통로.
	// InjectQueue에는 통로가 들어감. 데이터가 아님.
	// .mp3 파일이든, google tts든, 실시간 라디오 스트리밍 데이터든, 데이터가 흘러나올 통로.
	// 즉 InjectQueue에는 항상 "데이터 통로"가 쌓임.
	// 파일, HTTP 바디, TTS 스트림 등 뭐든 담을 수 있음.
	// 데이터가 준비될 때까지 reader.Read()가 블로킹하므로 실시간 스트림도 자연스럽게 처리됨
	reader io.Reader
	// 노래가 완료되면 done 채널로 메시지를 보내고 http connection을 끊음.
	done chan struct{}
}

var queueSize = 20
var InjectQueue = make(chan InjectItem, queueSize)

func InitInjectQueueWorker() {

	for item := range InjectQueue {

		// 오디오 청크를 저장할 버퍼(4096 바이트)를 생성.
		// 데이터 청크가 4096 바이트 만큼 쌓이면 UnBlocking 후 클라이언트에게 데이터를 전송함
		buf := make([]byte, 4096)

		for {

			// 데이터 통로에서 데이터 청크를 꺼내 buf에 채워준다.
			// 반환하는 n은 실제로 읽어들인 바이트 수
			// buf가 가득 찰 때까지 블로킹.
			n, err := item.reader.Read(buf)

			// 버퍼가 채워지면
			if n > 0 {
				// Broadcaster가 버퍼 내 오디오 청크를 서버와 연결된 모든 클라이언트에게 전송
				bc.Send(buf[:n])
			}

			// 오디오 청크를 모두 읽는 경우, 다음 데이터를 클라이언트에게 전송
			if err == io.EOF {
				break
			}
			// 오디오 청크를 읽는 도중 예외 발생시, 다음 데이터를 클라이언트 전송
			if err != nil {
				break
			}
		}
		// 읽기 완료 → 핸들러에게 커넥션 종료해도 된다고 신호
		close(item.done)
	}
}
