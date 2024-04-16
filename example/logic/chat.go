package logic

import (
	"bufio"
	"bytes"
	"github.com/swxctx/goai/baidu"
	"github.com/swxctx/goai/example/args"
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
	"io"
	"net/http"
	"time"
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

func chatBaidu(ctx *td.Context, arg *args.ChatDoArgsV1) (*args.ChatDoResultV1, *td.Rerror) {
	// 正常输出
	if !arg.Stream {
		resp, err := baidu.Chat(
			baidu.MODEL_FOR_35_8K,
			&baidu.ChatRequest{
				Messages: []baidu.MessageInfo{
					{
						Role:    baidu.CHAT_MESSAGE_ROLE_USER,
						Content: "你好",
					},
				},
			})
		if err != nil {
			xlog.Errorf("chatBaidu: err-> %v", err)
			return nil, td.RerrInternalServer.SetReason(err.Error())
		}
		return &args.ChatDoResultV1{
			Message: resp.Result,
		}, nil
	}

	// 流式输出
	resp, err := baidu.ChatStream(
		baidu.MODEL_FOR_35_8K,
		&baidu.ChatRequest{
			Messages: []baidu.MessageInfo{
				{
					Role:    baidu.CHAT_MESSAGE_ROLE_USER,
					Content: "你好",
				},
			},
		})
	if err != nil {
		xlog.Errorf("chatBaidu: err-> %v", err)
		return nil, td.RerrInternalServer.SetReason(err.Error())
	}

	// 读取流处理
	reader := bufio.NewReader(resp.Body)
	// 确保关闭
	defer resp.Body.Close()

	ctx.Stream(func(w io.Writer) bool {
		// 读取数据
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				// 记录错误，非EOF错误则终止循环
				xlog.Errorf("chatBaidu: err-> %v", err)
			}
			// 返回false通知Stream停止调用
			return false
		}

		trimMsg := bytes.TrimSpace(line)
		if len(trimMsg) == 0 {
			// 如果是空行，忽略，请求继续发送下一行
			return true
		}

		xlog.Infof("chatBaidu: line-> %s", string(line))

		// 写入一行数据到响应体
		w.Write(line)
		if flusher, ok := w.(http.Flusher); ok {
			// 确保数据发送到客户端
			flusher.Flush()
		}

		// 暂停一秒，方便观察效果
		time.Sleep(time.Duration(1) * time.Second)

		// 继续处理下一行
		return true
	})

	return new(args.ChatDoResultV1), nil
}
