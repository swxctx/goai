package xunfei

import (
	"fmt"

	"github.com/swxctx/xlog"
)

// ChatRequest 调用请求结构体
type ChatRequest struct {
	// 每个用户的id，用于区分不同用户[最大长度32]
	Uid string
	// 消息信息
	Messages []MessageInfo
	// 其他参数信息
	ChatParameter RequestParameterChat
}

// ChatStream 讯飞对话接口
func (c *Client) ChatStream(chatRequest *ChatRequest) (*StreamReader, error) {
	// 处理请求域
	requestUrl := fmt.Sprintf(c.baseUri, model_api_map[chatRequest.ChatParameter.Domain])

	// 处理请求参数
	requestParams := &xfRequest{
		Header: requestHeader{
			AppId: c.appId,
			Uid:   chatRequest.Uid,
		},
		Parameter: requestParameter{
			Chat: chatRequest.ChatParameter,
		},
		Payload: requestPayload{
			Message: requestPayloadMessage{
				Text: chatRequest.Messages,
			},
		},
	}

	// 处理请求链接
	finalUrl, err := assembleAuthUrl(requestUrl, c.apiKey, c.apiSecret)
	if err != nil {
		return nil, err
	}

	if c.debug {
		xlog.Debugf("xunfei: Chat request url-> %s", finalUrl)
	}

	// 发起连接
	conn, resp, err := c.webSocketDialer.Dial(finalUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("xunfei: Chat connect err-> %v", err)
	} else if resp.StatusCode != 101 {
		return nil, fmt.Errorf("xunfei: Chat connect status invalid, code-> %d", resp.StatusCode)
	}

	// 发送数据
	if err := conn.WriteJSON(requestParams); err != nil {
		return nil, fmt.Errorf("xunfei: Chat WriteJSON err-> %v", err)
	}

	return newStreamReader(conn), nil
}
