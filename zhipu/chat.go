package zhipu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/swxctx/ghttp"
	"github.com/swxctx/xlog"
)

// Chat 对话方法
func (c *Client) Chat(chatRequest *ChatRequest) (*ChatResponse, error) {
	if err := c.getAuthToken(); err != nil {
		return nil, err
	}

	chatRequest.Stream = false

	// new request
	req := ghttp.Request{
		Url:       c.baseUri,
		Method:    "POST",
		ShowDebug: c.debug,
		Body:      chatRequest,
	}
	req.AddHeader("Authorization", "Bearer "+c.authToken)
	req.AddHeader("Content-Type", "application/json")

	// send request
	resp, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("zhipu: Chat err, err is-> %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zhipu: Chat http response code not 200, code is -> %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// read body
	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("zhipu: Chat read resp body err-> %v", err)
	}
	if client.debug {
		xlog.Debugf("zhipu: chat resp-> %s", string(respBs))
	}

	// unmarshal data
	var (
		chatResp *ChatResponse
	)

	err = json.Unmarshal(respBs, &chatResp)
	if err != nil {
		return nil, fmt.Errorf("zhipu: Chat data unmarshal err-> %v", err)
	}
	return chatResp, nil
}

// ChatStream 流式对话方法
func (c *Client) ChatStream(chatRequest *ChatRequest) (*StreamReader, error) {
	if err := c.getAuthToken(); err != nil {
		return nil, err
	}

	chatRequest.Stream = true

	// new request
	req := ghttp.Request{
		Url:       c.baseUri,
		Method:    "POST",
		ShowDebug: c.debug,
		Body:      chatRequest,
	}
	req.AddHeader("Authorization", "Bearer "+c.authToken)
	req.AddHeader("Content-Type", "application/json")
	req.AddHeader("Connection", "keep-alive")

	// send request
	resp, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("zhipu: Chat err, err is-> %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zhipu: Chat http response code not 200, code is -> %d", resp.StatusCode)
	}

	// 交给外部调用逻辑处理
	return newStreamReader(resp, client.maxEmptyMessageCount), nil
}
