package git

import "context"

type Info struct {
	HtmlUrl string `json:"html_url"`
	Url     string `json:"url"`
	Number  int32  `json:"number"`
}

type Git interface {
	CreateFeature(ctx context.Context, repos, source, target string) error
	CreatePR(ctx context.Context, repos, source, target string) (*Info, error)
	GetPR(ctx context.Context, repos, source, target string) ([]*Info, error)
	Merge(ctx context.Context, repos string, num int32) error
	MergePR(ctx context.Context, repos string, source, target string) error
}
