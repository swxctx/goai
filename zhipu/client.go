package zhipu

import "github.com/swxctx/xlog"

/*
	接口文档：https://open.bigmodel.cn/dev/api#glm-4
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

	// auth token
	authToken string
	// 过期时间[存储的是过期时间节点的时间戳]
	expireIn int64
	// token过期时间，单位为：秒
	expSeconds int64

	// 是否调试模式[调试模式可以输出详细的信息]
	debug bool

	// 最大空消息数量
	maxEmptyMessageCount int
}

// NewClient 初始化请求客户端
func NewClient(apiKey string, debug ...bool) error {
	client = &Client{
		apiKey:               apiKey,
		baseUri:              "https://open.bigmodel.cn/api/paas/v4/chat/completions",
		maxEmptyMessageCount: 900,
		expSeconds:           3600,
	}
	if len(debug) > 0 {
		client.debug = debug[0]
	}

	if client.debug {
		xlog.SetLevel("debug")
	}

	// 初始化获取token
	return client.getAuthToken()
}

// GetAuthToken 返回auth token信息，比如在相同业务系统还需要用到这个Token
func GetAuthToken() (string, int64) {
	return client.authToken, client.expireIn
}

// RefreshAuthToken 刷新access token
func RefreshAuthToken() error {
	return client.refreshAuthToken()
}

// SetAuthTokenExp 设置auth token过期时间(单位为秒)
func SetAuthTokenExp(expSeconds int64) {
	client.expSeconds = expSeconds
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
