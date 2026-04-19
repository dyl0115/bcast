package main

import (
	"bytes"
	"io"
	"net/http"
)

func handleStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "audio/mpeg")
	ch := bc.Register()
	defer bc.Unregister(ch)

	for chunk := range ch {
		w.Write(chunk)
		w.(http.Flusher).Flush()
	}
}

func handleInject(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "read error", 500)
		return
	}

	InjectQueue <- bytes.NewReader(data) // bytes.Reader로 큐에 넣기
}
