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
	for {
		line, err := reader.ReadBytes('\n')
		// 结束
		if err == io.EOF {
			break
		}
		if err != nil {
			xlog.Errorf("chatBaidu: err-> %v", err)
			break
		}
		trimMsg := bytes.TrimSpace(line)
		if len(trimMsg) == 0 {
			continue
		}
		xlog.Infof("chatBaidu: line-> %s", string(line))
		ctx.ResponseWriter.Write(line)
		if flusher, ok := ctx.ResponseWriter.(http.Flusher); ok {
			flusher.Flush()
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
	// 确保关闭
	resp.Body.Close()

	return new(args.ChatDoResultV1), nil
}
