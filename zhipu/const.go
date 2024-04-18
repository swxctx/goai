package zhipu

const (
	// CHAT_MESSAGE_ROLE_USER 用户
	CHAT_MESSAGE_ROLE_USER string = "user"

	// CHAT_MESSAGE_ROLE_ASSISTANT 对话助手
	CHAT_MESSAGE_ROLE_ASSISTANT string = "assistant"

	// CHAT_MESSAGE_ROLE_SYSTEM 对话背景
	CHAT_MESSAGE_ROLE_SYSTEM string = "system"

	// CHAT_MESSAGE_ROLE_TOOL 工具调用
	CHAT_MESSAGE_ROLE_TOOL string = "tools"
)

const (
	// MODEL_VERSION_4
	MODEL_VERSION_4 = "glm-4"

	// MODEL_VERSION_3
	MODEL_VERSION_3 = "glm-3-turbo"
)

const (
	// TOOLS_TYPE_FOR_RETRIEVAL 知识库类型
	TOOLS_TYPE_FOR_RETRIEVAL = "retrieval"

	// TOOLS_TYPE_FOR_WEB_SEARCH 联网搜索
	TOOLS_TYPE_FOR_WEB_SEARCH = "web_search"
)
