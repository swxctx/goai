package aliyun

import "github.com/swxctx/xlog"

/*
	接口文档：https://help.aliyun.com/zh/dashscope/developer-reference/api-details?spm=a2c4g.11186623.0.0.e66d23edk4jpy6#b8ebf6b25eul6
*/

var (
	client *Client
)

// Client API请求客户端
type Client struct {
	// 基础请求api
	baseUri string

	// API Key
	apiKey string

	// 是否调试模式[调试模式可以输出详细的信息]
	debug bool

	// 最大空消息数量
	maxEmptyMessageCount int
}

// NewClient 初始化请求客户端
func NewClient(apiKey string, debug ...bool) {
	client = &Client{
		apiKey:               apiKey,
		baseUri:              "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation",
		maxEmptyMessageCount: 900,
	}
	if len(debug) > 0 {
		client.debug = debug[0]
	}

	if client.debug {
		xlog.SetLevel("debug")
	}

	return
}

// SetMaxEmptyMessageCount 最大空消息数量
func SetMaxEmptyMessageCount(count int) {
	client.maxEmptyMessageCount = count
}

// SetDebug debug开关
func SetDebug(debug bool) {
	client.debug = debug
}

// Chat 对话接口
func Chat(chatRequest *ChatRequest) (*ChatResponse, error) {
	return client.Chat(chatRequest)
}

// ChatStream 流式对话接口
func ChatStream(chatRequest *ChatRequest) (*StreamReader, error) {
	return client.ChatStream(chatRequest)
}
