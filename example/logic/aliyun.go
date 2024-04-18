package logic

import (
	"github.com/swxctx/goai/aliyun"
	"io"
	"net/http"

	"github.com/swxctx/goai/example/args"
	td "github.com/swxctx/malatd"
	"github.com/swxctx/xlog"
)

// chatAli
func chatAli(ctx *td.Context, arg *args.ChatDoArgsV1) (*args.ChatDoResultV1, *td.Rerror) {
	// 正常输出
	if !arg.Stream {
		resp, err := aliyun.Chat(&aliyun.ChatRequest{
			Model: aliyun.MODEL_VERSION_TURBO,
			Input: aliyun.Input{
				Prompt: arg.Content,
			},
			Parameters: aliyun.Parameters{
				ResultFormat: aliyun.RESULT_FORMAT_MESSAGE,
			},
		})
		if err != nil {
			xlog.Errorf("chatAli: err-> %v", err)
			return nil, td.RerrInternalServer.SetReason(err.Error())
		}
		return &args.ChatDoResultV1{
			Message: resp.Output.Choices[0].Message.Content,
		}, nil
	}

	// 流式输出
	return chatAliStream(ctx, arg)
}

// 流式输出(不接管流式数据处理)
func chatAliStream(ctx *td.Context, arg *args.ChatDoArgsV1) (*args.ChatDoResultV1, *td.Rerror) {
	streamReader, err := aliyun.ChatStream(&aliyun.ChatRequest{
		Model: aliyun.MODEL_VERSION_TURBO,
		Input: aliyun.Input{
			Prompt: arg.Content,
		},
		Parameters: aliyun.Parameters{
			ResultFormat:      aliyun.RESULT_FORMAT_MESSAGE,
			IncrementalOutput: true,
		},
	})
	if err != nil {
		xlog.Errorf("chatAliStream: err-> %v", err)
		return nil, td.RerrInternalServer.SetReason(err.Error())
	}

	defer streamReader.Close()
	ctx.Stream(func(w io.Writer) bool {
		data, err := streamReader.ReceiveFormat()
		if err != nil {
			xlog.Errorf("chatAliStream: Receive err-> %v", err)
			return false
		}
		if streamReader.IsFinish() {
			xlog.Infof("chatAliStream: receive finish...")
			return false
		}
		if streamReader.IsMaxEmptyLimit() {
			xlog.Infof("chatAliStream: empty line limit...")
			return false
		}
		if data == nil {
			xlog.Infof("chatAliStream: line is empty...")
			return true
		}

		xlog.Infof("resp: data-> %#v", data)

		// 写入一行数据到响应体
		w.Write(streamResponse(data.Output.Choices[0].Message.Content))
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
