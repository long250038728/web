package client

import (
	"context"
	_ "embed"
	"github.com/coze-dev/coze-go"
	"sync"
	"time"
)

//go:embed config/private_key.pem
var jwtOauthPrivateKey string

//go:embed config/client_id.txt
var jwtOauthClientID string

//go:embed config/public_key.txt
var jwtOauthPublicKeyID string

var oauth *coze.JWTOAuthClient

type Client struct {
	time  time.Time
	mutex sync.Mutex
}

func (c *Client) GetOAuth() (*coze.JWTOAuthClient, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	////超过900ttl,需要重新申请
	//if c.time.Add(time.Second*900).Unix() >= time.Now().Unix() {
	//	return oauth, nil
	//}

	// Default 15 minutes
	param := coze.NewJWTOAuthClientParam{
		ClientID:      jwtOauthClientID,
		PublicKey:     jwtOauthPublicKeyID,
		PrivateKeyPEM: jwtOauthPrivateKey,
	}
	oauth, err := coze.NewJWTOAuthClient(param, coze.WithAuthBaseURL(coze.CnBaseURL))
	if err != nil {
		return nil, err
	}
	c.time = time.Now()
	return oauth, nil
}

func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	oauth, err := c.GetOAuth()
	if err != nil {
		return "", err
	}

	req := &coze.GetJWTAccessTokenReq{}
	resp, err := oauth.GetAccessToken(ctx, req)
	return resp.AccessToken, err
}
