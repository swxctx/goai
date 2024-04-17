package xunfei

import "testing"

func TestChat(t *testing.T) {
	NewClient("appid", "apiKey", "apiSecret", true)
	// 调用Chat
	streamReader, err := ChatStream(&ChatRequest{
		Messages: []MessageInfo{
			{
				Role:    CHAT_MESSAGE_ROLE_USER,
				Content: "你好",
			},
		},
		ChatParameter: RequestParameterChat{
			Domain: DOMAIN_VERSION_30,
		},
	})
	if err != nil {
		t.Errorf("err-> %v", err)
	}

	defer streamReader.Close()
	for {
		data, err := streamReader.Receive()
		if err != nil {
			t.Errorf("err-> %v", err)
			break
		}
		if streamReader.IsFinish() {
			t.Logf("data load finish...")
			break
		}

		t.Logf("xunfei: resp-> %s", string(data))
	}
}

func TestChatFormat(t *testing.T) {
	NewClient("appid", "apiKey", "apiSecret", true)
	// 调用Chat
	streamReader, err := ChatStream(&ChatRequest{
		Messages: []MessageInfo{
			{
				Role:    CHAT_MESSAGE_ROLE_USER,
				Content: "你好",
			},
		},
		ChatParameter: RequestParameterChat{
			Domain: DOMAIN_VERSION_30,
		},
	})
	if err != nil {
		t.Errorf("err-> %v", err)
	}

	defer streamReader.Close()
	for {
		data, err := streamReader.ReceiveFormat()
		if err != nil {
			t.Errorf("err-> %v", err)
			break
		}
		if streamReader.IsFinish() {
			t.Logf("data load finish...")
			break
		}

		t.Logf("xunfei: resp-> %#v", data.Payload.Choices.Text[0].Content)
	}
}
