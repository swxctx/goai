package xunfei

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

type StreamReader struct {
	conn *websocket.Conn

	// 已经读取完了
	isFinish bool
	isEnd    bool
}

// newStreamReader
func newStreamReader(conn *websocket.Conn) *StreamReader {
	return &StreamReader{
		conn: conn,
	}
}

// Conn
func (streamReader *StreamReader) Conn() *websocket.Conn {
	return streamReader.conn
}

// IsFinish
func (streamReader *StreamReader) IsFinish() bool {
	return streamReader.isFinish
}

// Receive
func (streamReader *StreamReader) Receive() ([]byte, error) {
	if streamReader.isEnd {
		streamReader.isFinish = true
		return nil, nil
	}

	// 读取数据
	_, message, err := streamReader.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("xunfei: ChatStream ReadMessage err-> %v", err)
	}

	// 检查是否是最后一条了
	if strings.Contains(string(message), "\"status\":2,\"seq\":") {
		streamReader.isEnd = true
	}
	return message, nil
}

// ReceiveFormat
func (streamReader *StreamReader) ReceiveFormat() (*ChatResponse, error) {
	if streamReader.isEnd {
		streamReader.isFinish = true
		return nil, nil
	}

	var (
		chatResponse *ChatResponse
	)

	// 读取json数据
	if err := streamReader.conn.ReadJSON(&chatResponse); err != nil {
		return nil, fmt.Errorf("ReceiveFormat: ReadJSON err-> %v", err)
	}

	// 检查是否是最后一条了
	if chatResponse.Payload.Choices.Status == 2 {
		streamReader.isEnd = true
	}
	return chatResponse, nil
}

// Close
func (streamReader *StreamReader) Close() error {
	return streamReader.conn.Close()
}
