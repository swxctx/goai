package logic

import (
	"io"
	"net/http"

	"github.com/swxctx/goai/example/args"
	"github.com/swxctx/goai/zhipu"
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
)

// chatZP
func chatZP(ctx *td.Context, arg *args.ChatDoArgsV1) (*args.ChatDoResultV1, *td.Rerror) {
	// 正常输出
	if !arg.Stream {
		resp, err := zhipu.Chat(&zhipu.ChatRequest{
			Model: zhipu.MODEL_VERSION_3,
			Messages: []zhipu.MessageInfo{
				{
					Role:    zhipu.CHAT_MESSAGE_ROLE_USER,
					Content: arg.Content,
				},
			},
		})
		if err != nil {
			xlog.Errorf("chatZP: err-> %v", err)
			return nil, td.RerrInternalServer.SetReason(err.Error())
		}
		return &args.ChatDoResultV1{
			Message: resp.Choices[0].Message.Content,
		}, nil
	}

	// 流式输出
	return chatZPStream(ctx, arg)
}

// 流式输出(不接管流式数据处理)
func chatZPStream(ctx *td.Context, arg *args.ChatDoArgsV1) (*args.ChatDoResultV1, *td.Rerror) {
	streamReader, err := zhipu.ChatStream(&zhipu.ChatRequest{
		Model: zhipu.MODEL_VERSION_3,
		Messages: []zhipu.MessageInfo{
			{
				Role:    zhipu.CHAT_MESSAGE_ROLE_USER,
				Content: arg.Content,
			},
		},
	})
	if err != nil {
		xlog.Errorf("chatZPStream: err-> %v", err)
		return nil, td.RerrInternalServer.SetReason(err.Error())
	}

	defer streamReader.Close()
	ctx.Stream(func(w io.Writer) bool {
		data, err := streamReader.ReceiveFormat()
		if err != nil {
			xlog.Errorf("chatZPStream: Receive err-> %v", err)
			return false
		}
		if streamReader.IsFinish() {
			xlog.Infof("chatZPStream: receive finish...")
			return false
		}
		if streamReader.IsMaxEmptyLimit() {
			xlog.Infof("chatZPStream: empty line limit...")
			return false
		}
		if data == nil {
			xlog.Infof("chatZPStream: line is empty...")
			return true
		}

		xlog.Infof("resp: data-> %#v", data)

		// 写入一行数据到响应体
		w.Write(streamResponse(data.Choices[0].Message.Content))
		if flusher, ok := w.(http.Flusher); ok {
			// 确保数据发送到客户端
			flusher.Flush()
		}

		// 继续处理下一行
		return true
	})

	return &args.ChatDoResultV1{
		Message: STREAM_DONE_FLAG,
	}, nil
}
