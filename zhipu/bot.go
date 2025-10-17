package zhipu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/swxctx/ghttp"
	"github.com/swxctx/xlog"
)

// BotCreateConversation 创建智能体会话
func (c *Client) BotCreateConversation(appId string) (*BotCreateConversationResp, error) {
	if err := c.getAuthToken(); err != nil {
		return nil, err
	}

	// new request
	req := ghttp.Request{
		Url:       c.botUri + fmt.Sprintf("/v2/application/%s/conversation", appId),
		Method:    "POST",
		ShowDebug: c.debug,
	}
	req.AddHeader("Authorization", "Bearer "+c.authToken)
	req.AddHeader("Content-Type", "application/json")

	// send request
	resp, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("zhipu: CreateConversation err, err is-> %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zhipu: CreateConversation http response code not 200, code is -> %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// read body
	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("zhipu: CreateConversation read resp body err-> %v", err)
	}
	if client.debug {
		xlog.Debugf("zhipu: CreateConversation resp-> %s", string(respBs))
	}

	var (
		chatResp *BotCreateConversationResp
	)

	// unmarshal data
	err = json.Unmarshal(respBs, &chatResp)
	if err != nil {
		return nil, fmt.Errorf("zhipu: CreateConversation data unmarshal err-> %v", err)
	}
	return chatResp, nil
}

// BotChat 对话方法
func (c *Client) BotChat(chatRequest *BotChatRequest) (*BotChatResponse, error) {
	if err := c.getAuthToken(); err != nil {
		return nil, err
	}

	chatRequest.Stream = false

	// new request
	req := ghttp.Request{
		Url:       c.botUri + "/v3/application/invoke",
		Method:    "POST",
		ShowDebug: c.debug,
		Body:      chatRequest,
	}
	req.AddHeader("Authorization", "Bearer "+c.authToken)
	req.AddHeader("Content-Type", "application/json")

	// send request
	resp, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("zhipu: BotChatStream err, err is-> %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zhipu: BotChatStream http response code not 200, code is -> %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// read body
	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("zhipu: BotChatStream read resp body err-> %v", err)
	}
	if client.debug {
		xlog.Debugf("zhipu: BotChatStream resp-> %s", string(respBs))
	}

	var (
		chatResp *BotChatResponse
	)

	// unmarshal data
	err = json.Unmarshal(respBs, &chatResp)
	if err != nil {
		return nil, fmt.Errorf("zhipu: BotChatStream data unmarshal err-> %v", err)
	}
	return chatResp, nil
}

// BotChatStream 流式对话方法
func (c *Client) BotChatStream(chatRequest *BotChatRequest) (*StreamReader, error) {
	if err := c.getAuthToken(); err != nil {
		return nil, err
	}

	chatRequest.Stream = true

	// new request
	req := ghttp.Request{
		Url:       c.botUri + "/v3/application/invoke",
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
		return nil, fmt.Errorf("zhipu: BotChatStream err, err is-> %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zhipu: BotChatStream http response code not 200, code is -> %d", resp.StatusCode)
	}

	// 交给外部调用逻辑处理
	return newStreamReader(resp, client.maxEmptyMessageCount), nil
}
