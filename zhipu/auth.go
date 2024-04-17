package zhipu

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// getAuthToken 获取token
func (c *Client) getAuthToken() error {
	// 如果当前是有的，那么就使用当前的
	if len(c.authToken) > 0 && c.expireIn > (time.Now().Unix()+120)*1000 {
		return nil
	}

	return c.refreshAuthToken()
}

// refreshAuthToken 刷新token
func (c *Client) refreshAuthToken() error {
	parts := strings.Split(c.apiKey, ".")
	if len(parts) != 2 {
		return fmt.Errorf("zhipu: getAuthToken invalid apikey")
	}

	// 解析用户id及secret
	id, secret := parts[0], parts[1]

	// 过期时间
	expireIn := time.Now().Add(time.Second*time.Duration(c.expSeconds)).Unix() * 1000

	// Create the claims
	claims := jwt.MapClaims{
		"api_key":   id,
		"exp":       expireIn,
		"timestamp": time.Now().Unix() * 1000,
	}

	// Create a new JWT token with the secret as the signing key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["alg"] = "HS256"
	token.Header["sign_type"] = "SIGN"

	// Sign and get the complete encoded token as a string
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	c.authToken = signedToken
	c.expireIn = expireIn
	return nil
}
