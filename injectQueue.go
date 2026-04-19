package main

import "io"

// 오디오가 동시에 여러개 queue로 들어오면, 사용자가 뒤죽박죽 섞인 오디오를 듣게 됨.
// 요청된 여러 오디오를 순서대로 InjectQueue(고채널)에 넣어줌.
// 백그라운드로 돌아가는 고루틴이 InjectQueue에서 오디오 청크를 꺼내서 BroadCast.send()로 클라이언트에게 오디오를 스트리밍
// 최대 20개의 오디오까지 인큐 가능
var queueSize = 20
var InjectQueue = make(chan io.Reader, queueSize)

func InitInjectQueueWorker() {
	for reader := range InjectQueue {
		buf := make([]byte, 4096)
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				bc.Send(buf[:n])
			}
			if err == io.EOF {
				break
			}
		}
	}
}
