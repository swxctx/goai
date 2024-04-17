package zhipu

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/swxctx/ghttp"
	"io"
	"strings"
)

type StreamReader struct {
	// 请求response
	response *ghttp.Response
	// 读取响应
	reader *bufio.Reader

	// 已经读取完了
	isFinish bool
	// 最大空行超限
	maxEmptyLimit bool

	// 空消息数量
	emptyMsgCount int
	// 最大消息数量
	maxEmptyMessageCount int
}

// newStreamReader
func newStreamReader(response *ghttp.Response, maxEmptyMessageCount int) *StreamReader {
	return &StreamReader{
		response:             response,
		reader:               bufio.NewReader(response.Body),
		maxEmptyMessageCount: maxEmptyMessageCount,
	}
}

// Response
func (streamReader *StreamReader) Response() *ghttp.Response {
	return streamReader.response
}

// Reader
func (streamReader *StreamReader) Reader() *bufio.Reader {
	return streamReader.reader
}

// IsFinish
func (streamReader *StreamReader) IsFinish() bool {
	return streamReader.isFinish
}

// IsFinish
func (streamReader *StreamReader) IsMaxEmptyLimit() bool {
	return streamReader.maxEmptyLimit
}

// Receive
func (streamReader *StreamReader) Receive() ([]byte, error) {
	// 读取数据
	line, err := streamReader.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			streamReader.isFinish = true
			return nil, nil
		}
		return nil, fmt.Errorf("zhipu: ChatStream ReadBytes err-> %v", err)
	}

	// 数据处理
	parseLine := streamDataParse(line)
	if len(parseLine) == 0 {
		streamReader.emptyMsgCount++
		// 超过最大空数据限制
		if streamReader.emptyMsgCount > streamReader.maxEmptyMessageCount {
			streamReader.maxEmptyLimit = true
			return nil, nil
		}
		return []byte{}, nil
	}

	// 结束
	if string(parseLine) == " [DONE]\n" {
		streamReader.isFinish = true
		return nil, nil
	}

	return parseLine, nil
}

// ReceiveFormat
func (streamReader *StreamReader) ReceiveFormat() (*ChatResponse, error) {
	// 读取数据
	line, err := streamReader.reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			streamReader.isFinish = true
			return nil, nil
		}
		return nil, fmt.Errorf("zhipu: ChatStream ReadBytes err-> %v", err)
	}

	// 数据处理
	parseLine := streamDataParse(line)
	if len(parseLine) == 0 {
		streamReader.emptyMsgCount++
		// 超过最大空数据限制
		if streamReader.emptyMsgCount > streamReader.maxEmptyMessageCount {
			streamReader.maxEmptyLimit = true
			return nil, nil
		}
		return nil, nil
	}

	// 结束
	if string(parseLine) == " [DONE]\n" {
		streamReader.isFinish = true
		return nil, nil
	}

	var (
		chatResponse *ChatResponse
	)
	if err := json.Unmarshal(parseLine, &chatResponse); err != nil {
		return nil, fmt.Errorf("zhipu: ChatStream data unmarshal err-> %v", err)
	}
	return chatResponse, nil
}

// Close
func (streamReader *StreamReader) Close() {
	streamReader.response.Body.Close()
}

// streamDataParse 流式输出处理
func streamDataParse(line []byte) []byte {
	// 可能返回空格字符串
	trimMsg := bytes.TrimSpace(line)

	if len(trimMsg) == 0 {
		return []byte{}
	}

	// 接收处理数据
	trimmedLine := strings.TrimPrefix(string(trimMsg), "data:")

	return []byte(trimmedLine + "\n")
}
