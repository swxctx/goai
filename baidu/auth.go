package baidu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/swxctx/ghttp"
	"github.com/swxctx/xlog"
)

/**
Access Token 获取
接口文档：https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Ilkkrb0i5
*/

// accessTokenResponse 获取AccessToken响应信息
type accessTokenResponse struct {
	// 暂时用不到
	RefreshToken string `json:"refresh_token"`
	// 过期时间，单位为秒级(返回的是剩余的秒数)
	ExpiresIn int64 `json:"expires_in"`
	// token信息串
	AccessToken string `json:"access_token"`
	// 错误吗
	error string `json:"error"`
	// 错误描述信息
	ErrorDescription string `json:"error_description"`
}

// getAccessToken 获取api请求accessToken
func (c *Client) getAccessToken() error {
	// 先从缓存获取，如果有并且没有过期，那么直接使用就可以了
	if len(c.accessToken) > 0 && c.expireIn > time.Now().Unix() {
		return nil
	}

	return c.refreshAccessToken()
}

// refreshAccessToken 获取api请求accessToken
func (c *Client) refreshAccessToken() error {
	// new request
	req := ghttp.Request{
		Url:       fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", c.clientId, c.clientSecret),
		Method:    "GET",
		ShowDebug: c.debug,
	}

	// send request
	resp, err := req.Do()
	if err != nil {
		return fmt.Errorf("baidu: getAccessToken err, err is-> %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("baidu: getAccessToken http response code not 200, code is -> %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// read body
	respBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("baidu: getAccessToken read resp body err-> %v", err)
	}
	if client.debug {
		xlog.Debugf("resp: %s", string(respBs))
	}

	// unmarshal data
	var (
		tokenResp *accessTokenResponse
	)

	err = json.Unmarshal(respBs, &tokenResp)
	if err != nil {
		return fmt.Errorf("baidu: getAccessToken data unmarshal err-> %v", err)
	}
	if tokenResp == nil {
		return fmt.Errorf("baidu: getAccessToken data is nil")
	}

	// to client
	c.accessToken = tokenResp.AccessToken
	// 减少半小时，避免出现错误
	c.expireIn = tokenResp.ExpiresIn + time.Now().Unix() - 1800

	return nil
}
