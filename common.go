package goai

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/swxctx/ghttp"
	"io"
	"strings"
)

type (
	// StreamFunc 流式输出回调
	StreamFunc func(line []byte) bool
)

// StreamLogic 流式数据处理逻辑
func StreamLogic(resp *ghttp.Response, maxEmptyMessageCount int, streamFunc StreamFunc) error {
	// 读取数据
	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	// 空消息数量
	emptyMsgCount := 0

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("baidu: ChatStream ReadBytes err-> %v", err)
		}

		// 数据处理
		parseLine := StreamDataParse(line)
		if len(parseLine) == 0 {
			emptyMsgCount++
			if emptyMsgCount > maxEmptyMessageCount {
				return fmt.Errorf("baidu: ChatStream empty msg to long")
			}

			// 如果是空行，忽略，请求继续发送下一行
			continue
		}

		// 回调数据
		if !streamFunc(parseLine) {
			break
		}
	}

	return nil
}

// StreamDataParse 流式输出处理
func StreamDataParse(line []byte) []byte {
	// 可能返回空格字符串
	trimMsg := bytes.TrimSpace(line)

	if len(trimMsg) == 0 {
		return []byte{}
	}

	// 接收处理数据
	trimmedLine := strings.TrimPrefix(string(trimMsg), "data:")

	return []byte(trimmedLine + "\n")
}
