package xunfei

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// assembleAuthUrl 创建鉴权url apikey 即 hmac username
func assembleAuthUrl(hostUrl string, apiKey, apiSecret string) (string, error) {
	ul, err := url.Parse(hostUrl)
	if err != nil {
		return "", fmt.Errorf("xunfei: assembleAuthUrl Parse err-> %v", err)
	}

	// 签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	// date = "Tue, 28 May 2019 09:10:42 MST"
	// 参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}

	// 拼接签名字符串
	sign := strings.Join(signString, "\n")

	// 签名结果
	sha := hmacWithShaToBase64(sign, apiSecret)

	// 构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)

	// 将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	// 将编码后的字符串url encode后添加到url后面
	return hostUrl + "?" + v.Encode(), nil
}

// hmacWithShaToBase64
func hmacWithShaToBase64(data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}
