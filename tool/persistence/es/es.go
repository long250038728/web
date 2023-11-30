package es

import "github.com/olivere/elastic/v7"

//github.com/olivere/elastic/v7

func NewEs(config *Config) (*elastic.Client, error) {
	var options []elastic.ClientOptionFunc
	options = append(options, elastic.SetURL(config.Addr))
	options = append(options, elastic.SetBasicAuth(config.User, config.Password))
	options = append(options, elastic.SetSniff(false))

	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}
	return client, nil
}
