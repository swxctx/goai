package logic

import (
	"github.com/swxctx/goai/baidu"
	"github.com/swxctx/goai/example/args"
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
	"io"
	"net/http"
	"time"
)

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
	return chatBaiduStream(ctx, arg)
}

// 流式输出(不接管流式数据处理)
func chatBaiduStream(ctx *td.Context, arg *args.ChatDoArgsV1) (*args.ChatDoResultV1, *td.Rerror) {
	streamReader, err := baidu.ChatStream(
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

	defer streamReader.Close()
	ctx.Stream(func(w io.Writer) bool {
		line, err := streamReader.Receive()
		if err != nil {
			xlog.Errorf("chatBaidu: Receive err-> %v", err)
			return false
		}
		if streamReader.IsFinish() {
			xlog.Infof("chatBaidu: receive finish...")
			return false
		}
		if streamReader.IsMaxEmptyLimit() {
			xlog.Infof("chatBaidu: empty line limit...")
			return false
		}
		if len(line) == 0 {
			xlog.Infof("chatBaidu: line is empty...")
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
