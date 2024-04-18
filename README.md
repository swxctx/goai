# goai
Golang版本各厂商AI OpenAPI调用封装

- [通义千问](https://github.com/swxctx/goai/tree/main/aliyun)
- [文心一言](https://github.com/swxctx/goai/tree/main/baidu)
- [讯飞星火](https://github.com/swxctx/goai/tree/main/xunfei)
- [智谱AI](https://github.com/swxctx/goai/tree/main/zhipu)
- [使用示例Web程序](https://github.com/swxctx/goai/tree/main/example)
- [Android端调用example api示例](https://github.com/swxctx/aiStream)


#### 普通调用

```go
func TestChat(t *testing.T) {
	NewClient("apiKey", true)

	resp, err := Chat(&ChatRequest{
		Model: MODEL_VERSION_TURBO,
		Input: Input{
			Prompt: "你好",
		},
		Parameters: Parameters{
			ResultFormat: RESULT_FORMAT_MESSAGE,
		},
	})
	if err != nil {
		t.Errorf("err-> %v", err)
		return
	}
	t.Logf("ali: resp-> %#v", resp)
	t.Logf("ali: AI回复-> %s", resp.Output.Choices[0].Message.Content)
}
```

#### 流式输出调用

```go
func TestChatStream(t *testing.T) {
	NewClient("apiKey", true)

	streamReader, err := ChatStream(&ChatRequest{
		Model: MODEL_VERSION_TURBO,
		Input: Input{
			Prompt: "你好",
		},
		Parameters: Parameters{
			ResultFormat:      RESULT_FORMAT_MESSAGE,
			IncrementalOutput: true,
		},
	})
	if err != nil {
		t.Errorf("err-> %v", err)
		return
	}

	defer streamReader.Close()
	for {
		line, err := streamReader.Receive()
		if err != nil {
			t.Errorf("err-> %v", err)
			break
		}
		if streamReader.IsFinish() {
			t.Logf("read finish...")
			break
		}
		if streamReader.IsMaxEmptyLimit() {
			t.Errorf("read empty limit...")
			break
		}

		if len(line) == 0 {
			continue
		}
		t.Logf("ali: resp line-> %s, len-> %d", line, len(line))
	}
}
```

#### 与HTTP Web框架结合使用

```go
streamReader, err := aliyun.ChatStream(&aliyun.ChatRequest{
    Model: aliyun.MODEL_VERSION_TURBO,
    Input: aliyun.Input{
        Prompt: "你好",
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
    line, err := streamReader.Receive()
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
    if len(line) == 0 {
        xlog.Infof("chatAliStream: line is empty...")
        return true
    }

    xlog.Infof("chatAliStream: line-> %s", string(line))

    // 写入一行数据到响应体
    w.Write(line)
    if flusher, ok := w.(http.Flusher); ok {
        // 确保数据发送到客户端
        flusher.Flush()
    }

    // 继续处理下一行
    return true
})
```