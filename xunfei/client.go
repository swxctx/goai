package xunfei

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/swxctx/xlog"
)

/*
	接口文档
	https://www.xfyun.cn/doc/spark/Web.html#_1-%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E
*/

var (
	client *Client
)

// Client 百度API请求客户端
type Client struct {
	// 基础请求api
	baseUri string

	// appId 应用appid
	appId string
	// apiKey APIKey
	apiKey string
	// apiSecret APISecret
	apiSecret string

	// 是否调试模式[调试模式可以输出详细的信息]
	debug bool

	// 用于请求
	webSocketDialer *websocket.Dialer
}

// NewClient 初始化请求客户端
func NewClient(appId, apiKey, apiSecret string, debug ...bool) {
	client = &Client{
		appId:     appId,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseUri:   "wss://spark-api.xf-yun.com/%s/chat",
	}
	if len(debug) > 0 {
		client.debug = debug[0]
		xlog.SetLevel("debug")
	}

	// websocket
	client.webSocketDialer = &websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
}

// SetDebug debug开关
func SetDebug(debug bool) {
	client.debug = debug
}

// ChatStream 对话调用
func ChatStream(chatRequest *ChatRequest) (*StreamReader, error) {
	return client.ChatStream(chatRequest)
}
