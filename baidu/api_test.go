package baidu

import (
	"bufio"
	"bytes"
	"github.com/swxctx/xlog"
	"io"
	"testing"
)

func reloadClient() {
	if err := NewClient(
		"oRH73n8pqsdfsdfdf",
		"F45PIWIkwbcgquQsdfsdfsdfsd",
		true); err != nil {
		xlog.Errorf("NewClient: err-> %v", err)
	}
}

func TestAuth(t *testing.T) {
	reloadClient()
	accessToken, expireIn := GetAccessToken()
	t.Logf("GetAccessToken: token-> %s, expireIn-> %d", accessToken, expireIn)
}

func TestChat(t *testing.T) {
	reloadClient()
	resp, err := Chat(
		MODEL_FOR_35_8K,
		&ChatRequest{
			Messages: []MessageInfo{
				{
					Role:    CHAT_MESSAGE_ROLE_USER,
					Content: "你好",
				},
			},
		})
	if err != nil {
		t.Errorf("Chat: err-> %v", err)
		return
	}
	t.Logf("Chat: resp-> %s", resp.Result)
}

func TestChatStream(t *testing.T) {
	reloadClient()
	resp, err := ChatStream(
		MODEL_FOR_35_8K,
		&ChatRequest{
			Messages: []MessageInfo{
				{
					Role:    CHAT_MESSAGE_ROLE_USER,
					Content: "你好，你叫什么名字？",
				},
			},
		}, nil)
	if err != nil {
		t.Errorf("Chat: err-> %v", err)
		return
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
			t.Errorf("err-> %v", err)
			break
		}
		trimMsg := bytes.TrimSpace(line)
		if len(trimMsg) == 0 {
			continue
		}
		t.Logf("resp line-> %s, len-> %d", line, len(line))
	}
	// 确保关闭
	resp.Body.Close()
}
