// Command example is the malatd service project.
// The framework reference: https://github.com/swxctx/malatd
package __TPL__

// __API_TPL__ register PULL router
type __API_TPL__ interface {
	V1_Chat
}

type V1_Chat interface {
	Do(*ChatDoArgsV1) *ChatDoResultV1
}

type (
	ChatDoArgsV1 struct {
		// 厂商[1: 百度 2: 讯飞 3: 智谱]
		Platform int
		// 是否流式
		Stream bool
		// 对话内容
		Content string
	}
	ChatDoResultV1 struct {
		Message string
	}
)
