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
