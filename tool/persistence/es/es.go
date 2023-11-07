package es

import "github.com/olivere/elastic/v7"

//github.com/olivere/elastic/v7

type Persistence struct {
	Client *elastic.Client
}

func NewEs(addr string, user, password string) (*Persistence, error) {
	var options []elastic.ClientOptionFunc
	options = append(options, elastic.SetURL(addr))
	options = append(options, elastic.SetBasicAuth(user, password))
	options = append(options, elastic.SetSniff(false))

	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	return &Persistence{
		Client: client,
	}, nil
}
