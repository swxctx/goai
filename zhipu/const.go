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

// MessageContentType 定义消息内容类型常量
const (
	MessageContentTypeInput         = "input"          // 文本输入
	MessageContentTypeUploadFile    = "upload_file"    // 上传文件
	MessageContentTypeUploadAudio   = "upload_audio"   // 上传音频
	MessageContentTypeUploadImage   = "upload_image"   // 上传图片
	MessageContentTypeUploadVideo   = "upload_video"   // 上传视频
	MessageContentTypeSelectionList = "selection_list" // 可选项下拉框
)

const (
	MessageSceneTypeText         = "text"          // 默认，普通问答场景
	MessageSceneTypeComputerCall = "computer_call" // cogagent 节点回复场景
	MessageSceneTypeQA           = "qa"            // 问答节点回复场景
)
