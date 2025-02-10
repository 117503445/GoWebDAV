package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

// Middleware 是一个函数，它接受一个 http.Handler 并返回一个新的 http.Handler
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 捕获请求的完整内容
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to dump request")
		}

		// 创建一个自定义的 ResponseWriter 来捕获响应
		lrw := &loggingResponseWriter{ResponseWriter: w}

		// 调用下一个处理器
		next.ServeHTTP(lrw, r)

		// 捕获响应的完整内容
		responseDump := fmt.Sprintf(
			"HTTP/1.1 %d %s\n%s\n%s",
			lrw.statusCode,
			http.StatusText(lrw.statusCode),
			lrw.Header(),
			lrw.body.String(),
		)
		// fmt.Printf("---\nRequest:\n%s", string(requestDump))
		// fmt.Printf("Response:\n%s\n---\n", responseDump)
		d := fmt.Sprintf("---\n%s\n--\n%s\n---\n", string(requestDump), responseDump)

		fmt.Println(d)
		dirLogs := "./logs"
		if err := os.MkdirAll(dirLogs, 0755); err != nil {
			log.Warn().Err(err).Msg("Failed to create logs directory")
		} else {

			f := fmt.Sprintf("%v/%v.log", dirLogs, time.Now().Format("20060102.150405.000"))
			err = os.WriteFile(f, []byte(d), 0644)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to write file")
			}
		}
	})
}

// loggingResponseWriter 是一个自定义的 ResponseWriter，用于捕获响应状态码和响应体
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// WriteHeader 捕获状态码
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Write 捕获响应体
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b) // 捕获响应体
	return lrw.ResponseWriter.Write(b)
}
