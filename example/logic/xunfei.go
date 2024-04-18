package logic

import (
	"io"
	"net/http"

	"github.com/swxctx/goai/example/args"
	"github.com/swxctx/goai/xunfei"
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
)

// chatXF
func chatXF(ctx *td.Context, arg *args.ChatDoArgsV1) (*args.ChatDoResultV1, *td.Rerror) {
	// 正常输出
	streamReader, err := xunfei.ChatStream(&xunfei.ChatRequest{
		Messages: []xunfei.MessageInfo{
			{
				Role:    xunfei.CHAT_MESSAGE_ROLE_USER,
				Content: arg.Content,
			},
		},
		ChatParameter: xunfei.RequestParameterChat{
			Domain: xunfei.DOMAIN_VERSION_30,
		},
	})
	if err != nil {
		xlog.Errorf("chatXF: err-> %v", err)
		return nil, td.RerrInternalServer.SetReason(err.Error())
	}

	defer streamReader.Close()

	ctx.Stream(func(w io.Writer) bool {
		data, err := streamReader.ReceiveFormat()
		if err != nil {
			xlog.Errorf("chatXF: Receive err-> %v", err)
			return false
		}
		if streamReader.IsFinish() {
			xlog.Infof("chatXF: receive finish...")
			return false
		}
		if data == nil {
			xlog.Infof("chatXF: line is empty...")
			return true
		}

		xlog.Infof("resp: data-> %#v", data)

		// 写入一行数据到响应体
		w.Write(streamResponse(data.Payload.Choices.Content))
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
