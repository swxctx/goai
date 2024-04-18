package zhipu

// ChatRequest
type ChatRequest struct {
	// 所要调用的模型编码
	Model string `json:"model"`
	// 对话消息
	Messages []MessageInfo `json:"messages"`
	// 由用户端传参，需保证唯一性；用于区分每次请求的唯一标识，用户端不传时平台会默认生成。
	RequestId string `json:"request_id,omitempty"`
	// do_sample 为 true 时启用采样策略，do_sample 为 false 时采样策略 temperature、top_p 将不生效。默认值为 true。
	DoSample bool `json:"do_sample,omitempty"`
	// 是否使用流式输出
	Stream bool `json:"stream,omitempty"`
	// 随机性控制
	Temperature float32 `json:"temperature,omitempty"`
	// 核取样
	TopP float32 `json:"top_p,omitempty"`
	// 模型输出最大 tokens，最大输出为8192，默认值为1024
	MaxTokens int64 `json:"max_tokens,omitempty"`
	// 模型在遇到stop所制定的字符时将停止生成，目前仅支持单个停止词，格式为["stop_word1"]
	Stop []string `json:"stop,omitempty"`
	// 可供模型调用的工具列表,tools字段会计算 tokens ，同样受到tokens长度的限制
	Tools []ToolsInfo `json:"tools,omitempty"`
	// 终端用户的唯一ID，协助平台对终端用户的违规行为、生成违法及不良信息或其他滥用行为进行干预。ID长度要求：最少6个字符，最多128个字符。
	UserId string `json:"user_id,omitempty"`
}

// ToolsInfo
type ToolsInfo struct {
	/*
		工具类型
		retrieval: 知识库
		web_search: 联网搜索
	*/
	Type string `json:"type"`
	// 知识库使用
	Retrieval RetrievalInfo `json:"retrieval"`
}

// RetrievalInfo
type RetrievalInfo struct {
	// 知识库ID
	KnowledgeId string `json:"knowledge_id"`
	// 请求模型时的知识库模板
	PromptTemplate string `json:"prompt_template,omitempty"`
}

// WebSearch
type WebSearch struct {
	// 是否启用搜索，默认弃用
	Enable bool `json:"enable,omitempty"`
	// 强制搜索自定义关键内容，此时模型会根据自定义搜索关键内容返回的结果作为背景知识来回答用户发起的对话。
	SearchQuery string `json:"search_query"`
}

// MessageInfo 对话消息结构体
type MessageInfo struct {
	/*
		角色 取值为[system, user, assistant, tool]
		system: 用于设置对话背景
		user: 表示是用户的问题
		assistant: 表示AI的回复
		tool: 工具调用回复
	*/
	Role string `json:"role"`
	// 用户和AI的对话内容
	// 所有content的累计tokens需控制8192以内
	Content string `json:"content"`
}

// ChatResponse
type ChatResponse struct {
	// 任务ID
	Id string `json:"id"`
	// 请求创建时间，是以秒为单位的 Unix 时间戳
	Created int64 `json:"created"`
	// 模型名称
	Model string `json:"model"`
	// 请求ID
	RequestId string `json:"request_id"`
	// 当前对话的模型输出内容
	Choices []Choices `json:"choices"`
	// token消耗信息
	Usage Usage `json:"usage"`
}

type Choices struct {
	// 结果下标
	Index int `json:"index"`
	/*
		模型推理终止的原因。
		stop代表推理自然结束或触发停止词。
		tool_calls 代表模型命中函数。
		length代表到达 tokens 长度上限。
		sensitive 代表模型推理内容被安全审核接口拦截。请注意，针对此类内容，请用户自行判断并决定是否撤回已公开的内容。
		network_error 代表模型推理异常。
	*/
	FinishReason string `json:"finish_reason"`
	// 模型返回的文本信息
	Message MessageInfo `json:"message"`
	// 流式输出是使用这个字段
	Delta MessageInfo `json:"delta"`
}

// Usage
type Usage struct {
	// 用户输入的 tokens 数量
	PromptTokens int64 `json:"prompt_tokens"`
	// 模型输入的 tokens 数量
	CompletionTokens int64 `json:"completion_tokens"`
	// 总 tokens 数量
	TotalTokens int64 `json:"total_tokens"`
}
