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

type ChatRequest struct {
	// 聊天上下文信息
	Messages []MessageInfo `json:"messages"`
	// 随机性
	Temperature float32 `json:"temperature,omitempty"`
	// 多样性
	TopP float32 `json:"top_p,omitempty"`
	// 惩罚值
	PenaltyScore float32 `json:"penalty_score,omitempty"`

	// 流式输出[不用填]
	Stream bool `json:"stream,omitempty"`

	// 人设
	System string `json:"system,omitempty"`
	// 停止生成标识
	Stop []string `json:"stop,omitempty"`

	// 是否强制关闭实时搜索功能
	DisableSearch bool `json:"disableSearch,omitempty"`
	// 是否开启上角标返回
	EnableCitation bool `json:"enable_citation,omitempty"`
	// 是否开启搜索溯源信息
	EnableTrace bool `json:"enable_trace,omitempty"`

	// 指定最大输出token数
	MaxOutputTokens int64 `json:"max_output_tokens,omitempty"`
	// 表示最终用户唯一标识
	UserId string `json:"user_id"`
}

// MessageInfo 对话信息
type MessageInfo struct {
	/**
	角色
	user: 表示用户
	assistant: 表示对话助手
	function: 表示函数
	*/
	Role string `json:"role"`
	// 对话内容
	Content string `json:"content"`
}

type ChatResponse struct {
	// 本轮对话ID
	Id string `json:"id"`
	// 回包类型 chat.completion：多轮对话返回
	Object string `json:"object"`
	// 时间戳
	Created int64 `json:"created"`
	// 表示当前子句的序号。只有在流式接口模式下会返回该字段
	SentenceID int64 `json:"sentence_id"`
	// 表示当前子句是否是最后一句。只有在流式接口模式下会返回该字段
	IsEnd bool `json:"is_end"`
	// 当前生成的结果是否被截断
	IsTruncated bool `json:"is_truncated"`
	/**
	输出内容标识，说明：
	· normal：输出内容完全由大模型生成，未触发截断、替换
	· stop：输出结果命中入参stop中指定的字段后被截断
	· length：达到了最大的token数，根据EB返回结果is_truncated来截断
	· content_filter：输出内容被截断、兜底、替换为**等
	· function_call：调用了funtion call功能
	*/
	FinishReason string
	// 搜索结果
	SearchInfo *SearchInfo
	// 生成的结果
	Result string `json:"result"`
	// 表示用户输入是否存在安全风险，是否关闭当前会话，清理历史会话信息
	NeedClearHistory bool `json:"need_clear_history"`
	/**
	说明：
	· 0：正常返回
	· 其他：非正常
	*/
	Flag int `json:"flag"`
	// 当need_clear_history为true时，此字段会告知第几轮对话有敏感信息，如果是当前问题，ban_round=-1
	BanRound int `json:"ban_round"`
	// 消耗说明
	Usage *Usage `json:"usage"`
}

// SearchInfo 搜索结果信息
type SearchInfo struct {
	// 序号
	Index int64 `json:"index"`
	// 搜索结果URL
	Url string `json:"url"`
	// 搜索结果标题
	Title string `json:"title"`
}

// Usage token 消耗说明
type Usage struct {
	// 问题tokens数
	PromptTokens int64 `json:"prompt_tokens"`
	// 回答tokens数
	CompletionTokens int64 `json:"completion_tokens"`
	// tokens总数
	TotalTokens int64 `json:"total_tokens"`
}

// PluginUsage 插件消耗说明
type PluginUsage struct {
	// plugin名称，chatFile：chat file插件消耗的tokens
	Name string `json:"name"`
	// 解析文档tokens
	ParseTokens int64 `json:"parse_tokens"`
	// 摘要文档tokens
	AbstractTokens int64 `json:"abstract_tokens"`
	// 检索文档tokens
	SearchTokens int64 `json:"search_tokens"`
	// 总tokens
	TotalTokens int64 `json:"total_tokens"`
}

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
		xlog.Debugf("resp: %s", string(respBs))
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
func (c *Client) ChatStream(model string, chatRequest *ChatRequest) (*ghttp.Response, error) {
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

	return resp, nil
}
