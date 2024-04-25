package es

import "github.com/olivere/elastic/v7"

//github.com/olivere/elastic/v7

type ES struct {
	*elastic.Client
}

func NewEs(config *Config) (*ES, error) {
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(config.Address),
		elastic.SetBasicAuth(config.User, config.Password),
		elastic.SetSniff(false),
	}
	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}
	return &ES{Client: client}, nil
}
