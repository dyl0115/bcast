package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/spf13/cobra"
)

var injectCmd = &cobra.Command{
	Use:   "inject",
	Short: "스트리밍할 .mp3파일 주입합니다.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		//ffmpeg -i song.mp3 -f mp3 - | curl.exe -X POST --data-binary @- http://localhost:8336/inject

		filePath := args[0]

		// ffmpeg -i song.mp3 -f mp3 - 실행
		ffmpeg := exec.Command("ffmpeg", "-i", filePath, "-f", "mp3", "-")

		// ffmpeg stdout을 파이프로 가져옴
		ffmpegOut, err := ffmpeg.StdoutPipe()
		if err != nil {
			return fmt.Errorf("ffmpeg 파이프 생성 실패: %w", err)
		}

		if err := ffmpeg.Start(); err != nil {
			return fmt.Errorf("ffmpeg 시작 실패: %w", err)
		}

		// HTTP POST 요청 - ffmpeg stdout을 body로 스트리밍
		resp, err := http.Post("http://localhost:8336/inject", "audio/mpeg", ffmpegOut)
		if err != nil {
			ffmpeg.Process.Kill()
			return fmt.Errorf("HTTP POST 실패: %w", err)
		}
		defer resp.Body.Close()

		// ffmpeg 종료 대기
		if err := ffmpeg.Wait(); err != nil {
			return fmt.Errorf("ffmpeg 종료 오류: %w", err)
		}

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("서버 응답 [%d]: %s\n", resp.StatusCode, string(body))
		return nil
	},
}
