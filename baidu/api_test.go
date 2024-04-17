package baidu

import (
	"testing"

	"github.com/swxctx/xlog"
)

func reloadClient() {
	if err := NewClient(
		"apiKey",
		"secretKey",
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
	streamReader, err := ChatStream(
		MODEL_FOR_35_8K,
		&ChatRequest{
			Messages: []MessageInfo{
				{
					Role:    CHAT_MESSAGE_ROLE_USER,
					Content: "你好，你叫什么名字？",
				},
			},
		})
	if err != nil {
		t.Errorf("Chat: err-> %v", err)
		return
	}

	// 关闭
	defer streamReader.Close()

	// 读取流处理
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
			t.Errorf("read empty limir...")
			break
		}

		if len(line) == 0 {
			continue
		}
		t.Logf("resp line-> %s, len-> %d", line, len(line))
	}
}

func TestChatStreamFormat(t *testing.T) {
	reloadClient()
	streamReader, err := ChatStream(
		MODEL_FOR_35_8K,
		&ChatRequest{
			Messages: []MessageInfo{
				{
					Role:    CHAT_MESSAGE_ROLE_USER,
					Content: "你好，你叫什么名字？",
				},
			},
		})
	if err != nil {
		t.Errorf("Chat: err-> %v", err)
		return
	}

	// 关闭
	defer streamReader.Close()

	// 读取流处理
	for {
		chatResponse, err := streamReader.ReceiveFormat()
		if err != nil {
			t.Errorf("err-> %v", err)
			break
		}
		if streamReader.IsFinish() {
			t.Logf("read finish...")
			break
		}
		if streamReader.IsMaxEmptyLimit() {
			t.Errorf("read empty limir...")
			break
		}

		if chatResponse == nil {
			continue
		}
		t.Logf("resp line-> %#v", chatResponse)
	}
}
