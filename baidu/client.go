package baidu

import "github.com/swxctx/xlog"

var (
	client *Client
)

// Client 百度API请求客户端
type Client struct {
	// 基础请求api
	baseUri string

	// 应用的API Key
	clientId string
	// 应用的Secret Key
	clientSecret string

	// access token
	accessToken string
	// 过期时间[存储的是过期时间节点的时间戳]
	expireIn int64

	// 是否调试模式[调试模式可以输出详细的信息]
	debug bool

	// 最大空消息数量
	maxEmptyMessageCount int
}

// NewClient 初始化百度请求客户端
func NewClient(apiKey, secretKey string, debug ...bool) error {
	client = &Client{
		clientId:             apiKey,
		clientSecret:         secretKey,
		baseUri:              "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/%s?access_token=%s",
		maxEmptyMessageCount: 300,
	}
	if len(debug) > 0 {
		client.debug = debug[0]
	}

	if client.debug {
		xlog.SetLevel("debug")
	}

	// 初始化获取token
	return client.getAccessToken()
}

// GetAccessToken 返回access token信息，比如在相同业务系统还需要用到这个Token
func GetAccessToken() (string, int64) {
	return client.accessToken, client.expireIn
}

// RefreshAccessToken 刷新access token
func RefreshAccessToken() error {
	return client.refreshAccessToken()
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
func Chat(model string, chatRequest *ChatRequest) (*ChatResponse, error) {
	return client.Chat(model, chatRequest)
}

// ChatStream 流式对话接口
func ChatStream(model string, chatRequest *ChatRequest) (*StreamReader, error) {
	return client.ChatStream(model, chatRequest)
}
