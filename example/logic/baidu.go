package logic

import (
	"bufio"
	"github.com/swxctx/goai"
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

	// 流式输出(不接管流式数据处理)
	//return chatBaiduStream(ctx, arg)

	// 流式输出(接管流式数据处理)
	return chatBaiduStreamFunc(ctx, arg)
}

// 流式输出(不接管流式数据处理)
func chatBaiduStream(ctx *td.Context, arg *args.ChatDoArgsV1) (*args.ChatDoResultV1, *td.Rerror) {
	resp, err := baidu.ChatStream(
		baidu.MODEL_FOR_35_8K,
		&baidu.ChatRequest{
			Messages: []baidu.MessageInfo{
				{
					Role:    baidu.CHAT_MESSAGE_ROLE_USER,
					Content: "你好",
				},
			},
		}, nil)
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

		parseLine := goai.StreamDataParse(line)
		if len(parseLine) == 0 {
			// 如果是空行，忽略，请求继续发送下一行
			return true
		}
		xlog.Infof("chatBaidu: line-> %s", string(parseLine))

		// 写入一行数据到响应体
		w.Write(parseLine)
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

// 流式输出(接管流式数据处理)
func chatBaiduStreamFunc(ctx *td.Context, arg *args.ChatDoArgsV1) (*args.ChatDoResultV1, *td.Rerror) {
	_, err := baidu.ChatStream(
		baidu.MODEL_FOR_35_8K,
		&baidu.ChatRequest{
			Messages: []baidu.MessageInfo{
				{
					Role:    baidu.CHAT_MESSAGE_ROLE_USER,
					Content: "你好",
				},
			},
		}, func(line []byte) bool {
			// 返回给前端
			_, err := ctx.ResponseWriter.Write(line)
			if err != nil {
				// 处理写入错误
				return false
			}
			if flusher, ok := ctx.ResponseWriter.(http.Flusher); ok {
				flusher.Flush()
			}

			// 暂停一秒，方便观察效果
			time.Sleep(time.Duration(1) * time.Second)
			return true
		})
	if err != nil {
		xlog.Errorf("chatBaidu: err-> %v", err)
		return nil, td.RerrInternalServer.SetReason(err.Error())
	}

	return new(args.ChatDoResultV1), nil
}
