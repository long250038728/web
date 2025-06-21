package client

import (
	"context"
	_ "embed"
	"github.com/coze-dev/coze-go"
)

//go:embed config/private_key.pem
var jwtOauthPrivateKey string

//go:embed config/client_id.txt
var jwtOauthClientID string

//go:embed config/public_key.txt
var jwtOauthPublicKeyID string

type Client struct {
	oauth *coze.JWTOAuthClient
	auth  coze.Auth
}

func NewCozeClient() (CozeClientInterface, error) {
	param := coze.NewJWTOAuthClientParam{
		ClientID:      jwtOauthClientID,
		PublicKey:     jwtOauthPublicKeyID,
		PrivateKeyPEM: jwtOauthPrivateKey,
	}
	oauth, err := coze.NewJWTOAuthClient(param, coze.WithAuthBaseURL(coze.CnBaseURL))
	if err != nil {
		return nil, err
	}
	return &Client{oauth: oauth, auth: coze.NewJWTAuth(oauth, nil)}, nil
}

func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	req := &coze.GetJWTAccessTokenReq{}
	resp, err := c.oauth.GetAccessToken(ctx, req)
	return resp.AccessToken, err
}

func (c *Client) getApi(opts ...coze.CozeAPIOption) coze.CozeAPI {
	opts = append(opts, coze.WithBaseURL(coze.CnBaseURL))
	return coze.NewCozeAPI(c.auth, opts...)
}
