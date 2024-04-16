package logic

import (
	"github.com/swxctx/goai/example/args"
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
)

/**
常规调用
curl -X POST "http://127.0.0.1:8080/example/v1/chat/do" \
     -H "Connection: keep-alive" \
     -H "Content-Type: application/json" \
     -d '{
           "platform": 1,
           "stream": false,
           "content": "你好"
         }'

流式调用
curl -N -X POST "http://127.0.0.1:8080/example/v1/chat/do" \
     -H "Connection: keep-alive" \
     -H "Content-Type: application/json" \
     -d '{
           "platform": 1,
           "stream": true,
           "content": "你好"
         }'
*/

// Do handler
func V1_Chat_Do(ctx *td.Context, arg *args.ChatDoArgsV1) (*args.ChatDoResultV1, *td.Rerror) {
	switch arg.Platform {
	case 1:
		return chatBaidu(ctx, arg)
	default:
		xlog.Errorf("V1_Chat_Do: un support platform")
	}
	return new(args.ChatDoResultV1), nil
}
