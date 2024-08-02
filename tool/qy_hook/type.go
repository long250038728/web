package qy_hook

import "context"

type Hook interface {
	SendHook(ctx context.Context, content string, mobileList []string) error
}
