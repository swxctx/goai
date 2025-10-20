package zhipu

/*智能体相关参数*/

// BotCreateConversationResp 创建会话响应
type BotCreateConversationResp struct {
	Data struct {
		ConversationID string `json:"conversation_id"`
	} `json:"data"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// BotChatRequest 智能体对话
type BotChatRequest struct {
	AppID          string            `json:"app_id"`                     // 应用id（必传）
	ConversationID string            `json:"conversation_id,omitempty"`  // 会话id，未传则默认创建新的会话
	ThirdRequestID string            `json:"third_request_id,omitempty"` // 三方请求id
	Stream         bool              `json:"stream"`                     // 默认true，false时为同步调用
	Messages       []BotMessage      `json:"messages"`                   // 用户输入列表（必传）
	DocumentIDs    []string          `json:"document_ids,omitempty"`     // 问答类应用使用的文档ID集合
	KnowledgeIDs   []string          `json:"knowledge_ids,omitempty"`    // 问答类应用使用的知识ID集合
	SendLogEvent   bool              `json:"send_log_event,omitempty"`   // 是否实时推送过程日志
	NodeDocuments  []NodeDocumentMap `json:"nodeDocuments,omitempty"`    // 各节点对应的文档ID集合
}

// BotMessage 用户输入信息
type BotMessage struct {
	Role    string             `json:"role,omitempty"`    // user / assistant
	CallID  string             `json:"call_id,omitempty"` // 中断场景ID
	Type    string             `json:"type,omitempty"`    // text / computer_call / qa
	Content []BotMessageDetail `json:"content"`           // 具体内容
}

// BotMessageDetail 具体内容项
type BotMessageDetail struct {
	Type  string `json:"type"`          // input / upload_file / upload_audio / upload_image / upload_video / selection_list
	Value string `json:"value"`         // 用户输入或文件/图片URL等
	Key   string `json:"key,omitempty"` // 字段名称
}

// NodeDocumentMap 节点对应文档信息
type NodeDocumentMap struct {
	NodeID      int64   `json:"nodeId,omitempty"`      // 节点ID
	DocumentIDs []int64 `json:"documentIds,omitempty"` // 文档ID集合
}

// BotChatResponse 智能体接口返回结构
type BotChatResponse struct {
	ReqID          string     `json:"req_id"`              // 请求id
	ConversationID string     `json:"conversation_id"`     // 会话id
	AppID          string     `json:"app_id"`              // 智能体id
	Choices        []Choice   `json:"choices,omitempty"`   // 增量返回信息
	ErrorMsg       *ErrorCode `json:"error_msg,omitempty"` // 异常信息
}

// Choice 模型增量返回的单个结果
type Choice struct {
	Index        int             `json:"index"`                   // 结果下标
	FinishReason string          `json:"finish_reason,omitempty"` // stop/error
	Messages     *SyncMessages   `json:"messages,omitempty"`      // 同步调用结果
	Usage        []OpenUsageData `json:"usage,omitempty"`         // tokens统计
	Delta        *MessageDelta   `json:"delta,omitempty"`         // 当前会话输出消息体
}

// MessageDelta 模型增量消息
type MessageDelta struct {
	Content *MessageData `json:"content,omitempty"` // 模型推送消息
	Event   *EventData   `json:"event,omitempty"`   // 编排节点执行日志事件
}

// MessageData 模型输出内容（多类型）
type MessageData struct {
	Type     string      `json:"type"`                // text, image, video, all_tools, qa, computer_call
	Msg      interface{} `json:"msg,omitempty"`       // 根据type不同对应不同结构
	NodeID   string      `json:"node_id,omitempty"`   // 节点id
	NodeName string      `json:"node_name,omitempty"` // 节点名称
}

// EventData 节点执行日志事件
type EventData struct {
	NodeID    string       `json:"node_id,omitempty"`    // 节点id
	NodeName  string       `json:"node_name,omitempty"`  // 节点名称
	Type      string       `json:"type,omitempty"`       // node_processing / tool_processing / node_finish 等
	Content   *MessageData `json:"content,omitempty"`    // 输入输出内容
	Time      int64        `json:"time,omitempty"`       // 毫秒时间戳
	ToolCalls *ToolCalls   `json:"tool_calls,omitempty"` // 工具调用信息
}

// ToolCalls 工具调用封装
type ToolCalls struct {
	Type          string      `json:"type,omitempty"`            // function / retrieval / web_search
	ToolCallsData interface{} `json:"tool_calls_data,omitempty"` // 动态结构：函数、知识库、联网搜索
}

// SyncMessages 同步调用结果
type SyncMessages struct {
	Content *MessageData `json:"content,omitempty"` // 同步返回的内容
	Event   []EventData  `json:"event,omitempty"`   // 同步编排节点日志
}

// ErrorCode 错误信息
type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// OpenUsageData tokens统计信息
type OpenUsageData struct {
	Model            string `json:"model"`               // 模型名称
	NodeName         string `json:"node_name,omitempty"` // 节点名称
	InputTokenCount  int    `json:"input_token_count"`   // 输入tokens
	OutputTokenCount int    `json:"output_token_count"`  // 输出tokens
	TotalTokenCount  int    `json:"total_token_count"`   // 总tokens
}

//
// ------------------------- 以下为多场景消息体 -------------------------
//

// QaMsg 问答节点消息
type QaMsg struct {
	CallID     string   `json:"call_id,omitempty"`
	Question   string   `json:"question,omitempty"`
	AnswerType string   `json:"answer_type,omitempty"` // option / input
	Options    []string `json:"options,omitempty"`
}

// ComputerCallMsg cogagent节点消息
type ComputerCallMsg struct {
	CallID    string `json:"call_id,omitempty"`
	Task      string `json:"task,omitempty"`
	Action    string `json:"action,omitempty"`
	Status    string `json:"status,omitempty"`
	Plan      string `json:"plan,omitempty"`
	Operation string `json:"operation,omitempty"`
	Level     string `json:"level,omitempty"` // critical / sensitive / general / end
}

// AllToolsMsg all_tools类型消息
type AllToolsMsg struct {
	Code string `json:"code,omitempty"`
	File string `json:"file,omitempty"`
	Text string `json:"text,omitempty"`
}

// VideoOrImageMsg 图片或视频生成消息
type VideoOrImageMsg struct {
	URL      string `json:"url,omitempty"`
	CoverURL string `json:"cover_url,omitempty"`
}

// FunToolCallsData 函数调用日志
type FunToolCallsData struct {
	ActionKey string `json:"action_key,omitempty"`
	Params    string `json:"params,omitempty"`
	Output    string `json:"output,omitempty"`
}

// KnowToolCallsData 知识库调用日志
type KnowToolCallsData struct {
	Input     string `json:"input,omitempty"`
	SliceInfo string `json:"slice_info,omitempty"`
}

// OpenWebSearchData 联网搜索日志
type OpenWebSearchData struct {
	Input   string `json:"input,omitempty"`
	Refer   string `json:"refer,omitempty"`
	Title   string `json:"title,omitempty"`
	Link    string `json:"link,omitempty"`
	Content string `json:"content,omitempty"`
	Media   string `json:"media,omitempty"`
	Icon    string `json:"icon,omitempty"`
}

// LoopActionData 循环节点轮次事件
type LoopActionData struct {
	CurrentRound int    `json:"current_round"`
	TotalRound   int    `json:"total_round"`
	Status       string `json:"status"`        // processing / finished
	FinishReason string `json:"finish_reason"` // 执行完成原因
}
