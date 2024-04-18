package aliyun

import "testing"

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

func TestChatStreamFormat(t *testing.T) {
	NewClient("sk-f76f969c31f04768ae718074c77a50d9", true)

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
		t.Logf("ali: resp data-> %#v", data)
	}
}
