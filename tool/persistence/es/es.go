package es

import "github.com/olivere/elastic/v7"

//github.com/olivere/elastic/v7

func NewEs(addr string, user, password string) (*elastic.Client, error) {
	var options []elastic.ClientOptionFunc
	options = append(options, elastic.SetURL(addr))
	options = append(options, elastic.SetBasicAuth(user, password))
	options = append(options, elastic.SetSniff(false))

	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}
	return client, nil
}
