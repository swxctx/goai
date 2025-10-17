package zhipu

import "testing"

func TestChat(t *testing.T) {
	NewClient("apiKey", true)

	resp, err := Chat(&ChatRequest{
		Model: MODEL_VERSION_3,
		Messages: []MessageInfo{
			{
				Role:    CHAT_MESSAGE_ROLE_USER,
				Content: "你好",
			},
		},
	})
	if err != nil {
		t.Errorf("err-> %v", err)
		return
	}
	t.Logf("zhipu: resp-> %#v", resp)
	t.Logf("zhipu: AI回复-> %s", resp.Choices[0].Message.Content)
}

func TestChatStream(t *testing.T) {
	NewClient("apiKey", true)

	streamReader, err := ChatStream(&ChatRequest{
		Model: MODEL_VERSION_4,
		Messages: []MessageInfo{
			{
				Role:    CHAT_MESSAGE_ROLE_USER,
				Content: "你好",
			},
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
		t.Logf("zhipu: resp line-> %s, len-> %d", line, len(line))
	}
}

func TestChatStreamFormat(t *testing.T) {
	NewClient("apiKey", true)

	streamReader, err := ChatStream(&ChatRequest{
		Model: MODEL_VERSION_3,
		Messages: []MessageInfo{
			{
				Role:    CHAT_MESSAGE_ROLE_USER,
				Content: "你好",
			},
		},
	})
	if err != nil {
		t.Errorf("err-> %v", err)
		return
	}

	defer streamReader.Close()
	for {
		data, err := streamReader.ReceiveFormat()
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

		if data == nil {
			continue
		}
		t.Logf("zhipu: resp line-> %#v", data)
		t.Logf("zhipu: resp 回复-> %s", data.Choices[0].Delta.Content)
	}
}

// go test -v -run=TestBotCreateConversation
func TestBotCreateConversation(t *testing.T) {
	NewClient("apikey", true)
	resp, err := BotCreateConversation("智能体appid")
	if err != nil {
		t.Errorf("err %v", err)
		return
	}
	t.Logf("cid-> %s", resp.Data.ConversationID)
}

// go test -v -run TestBotChat
func TestBotChat(t *testing.T) {
	NewClient("apikey", true)

	resp, err := BotChat(&BotChatRequest{
		AppID:          "智能体appid",
		ConversationID: "创建的对话id",
		Messages: []BotMessage{
			{
				Role: CHAT_MESSAGE_ROLE_USER,
				Content: []BotMessageDetail{
					{
						Type:  MessageContentTypeInput,
						Value: "你好",
					},
				},
			},
		},
	})
	if err != nil {
		t.Errorf("err-> %v", err)
		return
	}
	t.Logf("zhipu: resp-> %#v", resp)
	t.Logf("zhipu: AI回复-> %v", resp.Choices[0].Delta.Content.Msg.(string))
}

// go test -v -run=TestBotChatStream
func TestBotChatStream(t *testing.T) {
	NewClient("apikey", true)

	streamReader, err := BotChatStream(&BotChatRequest{
		AppID:          "智能体appid",
		ConversationID: "创建的对话id",
		Messages: []BotMessage{
			{
				Role: CHAT_MESSAGE_ROLE_USER,
				Content: []BotMessageDetail{
					{
						Type:  MessageContentTypeInput,
						Value: "你好",
					},
				},
			},
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
		t.Logf("zhipu: resp line-> %s, len-> %d", line, len(line))
	}
}
