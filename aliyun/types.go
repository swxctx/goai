package aliyun

// ChatRequest
type ChatRequest struct {
	/*
		指定用于对话的通义千问模型名
		目前可选择：
		qwen-turbo
		qwen-plus
		qwen-max
		qwen-max-0403
		qwen-max-0107
		qwen-max-1201
		qwen-max-longcontext
	*/
	Model string `json:"model"`
	// 输入参数
	Input Input `json:"input,omitempty"`
	// 模型参数设置
	Parameters Parameters `json:"parameters,omitempty"`
}

// Input
type Input struct {
	// 提示词
	Prompt string `json:"prompt"`
	// 历史消息列表
	Messages []MessageInfo `json:"messages,omitempty"`
}

// MessageInfo
type MessageInfo struct {
	// system、user、assistant和tool
	Role string `json:"role,omitempty"`
	// 消息内容
	Content string `json:"content,omitempty"`
}

// Parameters
type Parameters struct {
	// "text"表示旧版本的text
	// "message"表示兼容openai的message
	ResultFormat string `json:"result_format,omitempty"`
	// 随机种子
	Seed int64 `json:"seed,omitempty"`
	// 输出最大token限制
	MaxTokens int64 `json:"max_tokens,omitempty"`
	// 概率阀值
	TopP float32 `json:"top_p,omitempty"`
	// 候选集大小
	TopK int `json:"top_k,omitempty"`
	// 模型重复度控制
	RepetitionPenalty float32 `json:"repetition_penalty,omitempty"`
	// 多样性程度控制
	Temperature float32 `json:"temperature,omitempty"`
	// 停止生成标识
	Stop []string `json:"stop,omitempty"`
	// 是否启用联网搜索
	EnableSearch bool `json:"enable_search,omitempty"`
	// 控制流输出模式，如果设置为true，那么采用的就是增量输出
	IncrementalOutput bool `json:"incremental_output,omitempty"`
}

// ChatResponse
type ChatResponse struct {
	// 本次输出
	Output Output `json:"output"`
	// 请求ID
	RequestId string `json:"request_id"`
}

// Output
type Output struct {
	// 本次算法输出内容
	Text string `json:"text"`
	// 有三种情况：正在生成时为null，生成结束时如果由于停止token导致则为stop，生成结束时如果因为生成长度过长导致则为length。
	FinishReason string `json:"finish_reason"`
	// 入参result_format=message时候的返回值
	Choices []Choices `json:"choices"`
}

// Choices
type Choices struct {
	/*
		停止原因，null：生成过程中
		stop：stop token导致结束
		length：生成长度导致结束
	*/
	FinishReason string `json:"finish_reason"`
	// 消息数据
	Message MessageInfo `json:"message"`
}

// Usage
type Usage struct {
	// 本次请求算法输出内容的 token 数目。
	OutputTokens int64 `json:"output_tokens"`
	// 本次请求输入内容的 token 数目。在打开了搜索的情况下，输入的 token 数目因为还需要添加搜索相关内容支持，所以会超出客户在请求中的输入。
	InputTokens int64 `json:"input_tokens"`
}
