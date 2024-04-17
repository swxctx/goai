package xunfei

/*
# 参数构造示例如下
{
        "header": {
            "app_id": "12345",
            "uid": "12345"
        },
        "parameter": {
            "chat": {
                "domain": "generalv3.5",
                "temperature": 0.5,
                "max_tokens": 1024,
            }
        },
        "payload": {
            "message": {
                # 如果想获取结合上下文的回答，需要开发者每次将历史问答信息一起传给服务端，如下示例
                # 注意：text里面的所有content内容加一起的tokens需要控制在8192以内，开发者如有较长对话需求，需要适当裁剪历史信息
                "text": [
                    {"role":"system","content":"你现在扮演李白，你豪情万丈，狂放不羁；接下来请用李白的口吻和用户对话。"} #设置对话背景或者模型角色
                    {"role": "user", "content": "你是谁"} # 用户的历史问题
                    {"role": "assistant", "content": "....."}  # AI的历史回答结果
                    # ....... 省略的历史对话
                    {"role": "user", "content": "你会做什么"}  # 最新的一条问题，如无需上下文，可只传最新一条问题
                ]
        }
    }
}
*/

// xfRequest 请求参数
type xfRequest struct {
	// header
	Header requestHeader `json:"header"`
	// params
	Parameter requestParameter `json:"parameter"`
	// 消息体
	Payload requestPayload `json:"payload"`
}

// requestHeader header部分
type requestHeader struct {
	// 应用appid，从开放平台控制台创建的应用中获取
	AppId string `json:"app_id"`
	// 每个用户的id，用于区分不同用户[最大长度32]
	Uid string `json:"uid,omitempty"`
}

// requestParameter parameter
type requestParameter struct {
	// 详细参数
	Chat RequestParameterChat `json:"chat"`
}

// RequestParameterChat parameter.chat部分
type RequestParameterChat struct {
	/*
		指定访问的领域:
		general 指向V1.5版本;
		generalv2 指向V2版本;
		generalv3 指向V3版本;
		generalv3.5 指向V3.5版本;
	*/
	Domain string `json:"domain"`
	// 采样阀值 取值范围 (0，1] ，默认值0.5
	Temperature float64 `json:"temperature,omitempty"`
	/*
		模型回答的token最大长度
		V1.5取值为[1,4096]
		V2.0、V3.0和V3.5取值为[1,8192]，默认为2048
	*/
	MaxTokens int `json:"max_tokens,omitempty"`
	// 非等概率 取值为[1，6],默认为4
	TopK int `json:"top_k,omitempty"`
	// 用于关联用户会话[需要保障用户下的唯一性]
	ChatId string `json:"chat_id,omitempty"`
}

// requestPayload
type requestPayload struct {
	// 消息信息
	Message requestPayloadMessage `json:"message"`
}

// requestPayloadMessage
type requestPayloadMessage struct {
	Text []MessageInfo `json:"text"`
}

// MessageInfo 对话消息结构体
type MessageInfo struct {
	// 角色 取值为[system,user,assistant]
	// system用于设置对话背景，user表示是用户的问题，assistant表示AI的回复
	Role string `json:"role"`
	// 用户和AI的对话内容
	// 所有content的累计tokens需控制8192以内
	Content string `json:"content"`
}

// ChatResponse 请求响应
type ChatResponse struct {
	Header  ResponseHeader  `json:"header"`
	Payload ResponsePayload `json:"payload"`
}

// ResponseHeader
type ResponseHeader struct {
	// 错误码，0表示正常，非0表示出错
	Code int `json:"code"`
	// 会话是否成功的描述信息
	Message string `json:"message"`
	// 会话的唯一id，用于讯飞技术人员查询服务端会话日志使用,出现调用错误时建议留存该字段
	Sid string `json:"sid"`
	/*
		文本响应状态，取值为[0,1,2]
		0: 代表首个文本结果
		1: 代表中间文本结果
		2: 代表最后一个文本结果
	*/
	Status int `json:"status"`
}

// ResponsePayload
type ResponsePayload struct {
	// 回复信息
	Choices ChoicesPayload `json:"choices"`
	// token消耗信息
	Usage UsagePayload `json:"usage"`
}

// ChoicesPayload
type ChoicesPayload struct {
	/*
		文本响应状态，取值为[0,1,2]
		0: 代表首个文本结果
		1: 代表中间文本结果
		2: 代表最后一个文本结果
	*/
	Status int `json:"status"`
	// 返回的数据序号，取值为[0,9999999]
	Seq int `json:"seq"`
	// 回复的内容
	Text []ChoiceText `json:"text"`
	// 兼容处理
	MessageInfo
}

// ChoiceText
type ChoiceText struct {
	// AI的回答内容
	Content string `json:"content"`
	// 角色标识，固定为assistant，标识角色为AI
	Role string `json:"role"`
	// 结果序号，取值为[0,10]; 当前为保留字段，开发者可忽略
	Index int `json:"index"`
}

// UsagePayload
type UsagePayload struct {
	Text Usage `json:"text"`
}

// Usage
type Usage struct {
	// 提问消耗tokens
	QuestionTokens int64 `json:"question_tokens"`
	// 包含历史问题的总tokens大小
	PromptTokens int64 `json:"prompt_tokens"`
	// 回答的tokens大小
	CompletionTokens int64 `json:"completion_tokens"`
	// prompt_tokens和completion_tokens的和，也是本次交互计费的tokens大小
	TotalTokens int64 `json:"total_tokens"`
}
