package baidu

import (
	"encoding/json"
	"fmt"
	"github.com/swxctx/ghttp"
	"github.com/swxctx/xlog"
	"io/ioutil"
	"net/http"
)

/**
对话相关api
接口文档：https://cloud.baidu.com/doc/WENXINWORKSHOP/s/jlil56u11
*/

// Chat 对话方法
func (c *Client) Chat(model string, chatRequest *ChatRequest) (*ChatResponse, error) {
	if err := c.getAccessToken(); err != nil {
		return nil, err
	}

	chatRequest.Stream = false

	// new request
	req := ghttp.Request{
		Url:       fmt.Sprintf(c.baseUri, model, c.accessToken),
		Method:    "POST",
		ShowDebug: c.debug,
		Body:      chatRequest,
	}

	// send request
	resp, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("baidu: Chat err, err is-> %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("baidu: Chat http response code not 200, code is -> %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// read body
	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("baidu: Chat read resp body err-> %v", err)
	}
	if client.debug {
		xlog.Debugf("baidu: chat resp-> %s", string(respBs))
	}

	// unmarshal data
	var (
		chatResp *ChatResponse
	)

	err = json.Unmarshal(respBs, &chatResp)
	if err != nil {
		return nil, fmt.Errorf("baidu: Chat data unmarshal err-> %v", err)
	}
	return chatResp, nil
}

// ChatStream 流式对话方法
func (c *Client) ChatStream(model string, chatRequest *ChatRequest) (*StreamReader, error) {
	if err := c.getAccessToken(); err != nil {
		return nil, err
	}

	chatRequest.Stream = true

	// new request
	req := ghttp.Request{
		Url:       fmt.Sprintf(c.baseUri, model, c.accessToken),
		Method:    "POST",
		ShowDebug: c.debug,
		Body:      chatRequest,
	}
	req.AddHeader("Connection", "keep-alive")

	// send request
	resp, err := req.Do()
	if err != nil {
		return nil, fmt.Errorf("baidu: ChatStream err, err is-> %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("baidu: ChatStream http response code not 200, code is -> %d", resp.StatusCode)
	}

	// 交给外部调用逻辑处理
	return newStreamReader(resp, client.maxEmptyMessageCount), nil
}
