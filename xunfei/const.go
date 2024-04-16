package xunfei

const (
	// CHAT_MESSAGE_ROLE_USER 用户
	CHAT_MESSAGE_ROLE_USER string = "user"

	// CHAT_MESSAGE_ROLE_ASSISTANT 对话助手
	CHAT_MESSAGE_ROLE_ASSISTANT string = "assistant"

	// CHAT_MESSAGE_ROLE_SYSTEM 对话背景
	CHAT_MESSAGE_ROLE_SYSTEM string = "system"
)

const (
	// 1.5版本
	DOMAIN_VERSION_15 = "general"

	// 2.0版本
	DOMAIN_VERSION_20 = "generalv2"

	// 3.0版本
	DOMAIN_VERSION_30 = "generalv3"

	// 3.5版本
	DOMAIN_VERSION_35 = "generalv3.5"
)

var (
	// 域对应关系
	model_api_map = map[string]string{
		DOMAIN_VERSION_15: "v1.1",
		DOMAIN_VERSION_20: "v2.1",
		DOMAIN_VERSION_30: "v3.1",
		DOMAIN_VERSION_35: "v3.5",
	}
)
